package puppetmaster

import (
	"log"
)

func Example() {
	client, err := NewClient("https://puppet-master.io/api/v1/teams/my-team", "theapitokenigot")
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	client.EnableDebugLogs()

	jobs, err := client.GetJobs(1, 100)
	if err != nil {
		log.Fatalf("failed to fetch jobs: %v", err)
	}

	for _, job := range jobs.Jobs {
		log.Printf("Job ID %v", job.UUID)
	}

	newJob := &JobRequest{
		Code: `
import {getIp} from 'shared';

await page.goto(vars.page);
const ip = await getIp(page);

logger.info(ip);
results.ip = ip;
`,
		Modules: map[string]string{
			"shared": `
export async function getIp(page) {
  const text = await page.evaluate(() => document.querySelector('body').textContent);
  return text.split(":")[1];
}
`,
		},
		Vars: map[string]string{
			"page": "http://ifcfg.co",
		},
	}

	createdJob, err := client.CreateJob(newJob)
	if err != nil {
		log.Fatalf("failed to create job: %v", err)
	}

	log.Printf("Created job %v", createdJob.UUID)

	retrievedJob, err := client.GetJob(createdJob.UUID)
	if err != nil {
		log.Fatalf("failed to retrieve single job: %v", err)
	}

	log.Printf("Retrieved job: %v", retrievedJob.UUID)

	err = client.DeleteJob(retrievedJob.UUID)
	if err != nil {
		log.Fatalf("failed to delete job: %v", err)
	}

	log.Printf("Done executing lifecycle of job %v", retrievedJob.UUID)
}
