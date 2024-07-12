# **Jenkins Webhook Setup Tutorial for Separated Repositories**

Welcome to the Jenkins Webhook Setup Tutorial! In this guide, we will walk you through the process of configuring Jenkins webhooks to support separated repositories, ensuring accurate and meaningful DORA (DevOps Research and Assessment) metrics.

---

## **Unveiling the Challenge**

Traditionally, Jenkins assumes a direct alignment between the repository cloned during a job and the business logic repository essential for DORA metrics. However, real-world complexities emerge in scenarios where Jenkins jobs may originate from a repository distinct from the one housing the critical business logic.

Consider a nuanced example featuring two pivotal Git repositories: "webapp" and "jenkinsfiles." Within the DevLake project, the "webapp" is intricately linked to the GitHub repository "webapp" (say, repo1). In parallel, the Jenkins connection is meticulously configured to interface with the repository "jenkinsfiles" (say, repo2). Herein lies the challenge: during a Jenkins deployment job, the Git SHA is extracted from repo2 ("jenkinsfiles") rather than the anticipated repo1 ("webapp"). Intriguingly, the sole connection Jenkins has back to repo1 is orchestrated through a Jenkins input parameter coined "TAG," serving as a symbolic representation of a Git ref within repo1.

In essence, the conventional approach of assuming a direct correlation between the Jenkins job's cloned repository and the business logic repository faces a substantial limitation. This discrepancy becomes particularly pronounced when navigating the intricacies of deployment scenarios involving a "devops-setting commit" in repo2 (jenkinsfiles) and a simultaneous "business-code commit" in both repo1 (webapp). To achieve the coveted precision in measuring DORA metrics, we must pivot our attention to the commit originating from the business-code repository.

**Important Note:** Combining GitHub and Jenkins connections in the same project is insufficient for this use case. A deployment involves a "devops-setting commit" in repo2 (jenkinsfiles) and a "business-code commit" in both repo1 (webapp). To measure DORA metrics accurately, we must use the commit from the business-code repository.

*At this juncture, we arrive at the crux of our tutorial.*

---

## **Setting Up Jenkins Webhook for Separated Repositories:**

### **Step 1: Configure GitHub Connection**

1. Open your Jenkins instance and navigate to the DevLake project for the "webapp."
2. In the project settings, configure the GitHub connection to point to repo1 ("webapp").

### **Step 2: Add Jenkins Connection**

1. Still in the DevLake project, add a Jenkins connection pointing to repo2 ("jenkinsfiles").
2. Ensure that the Jenkins input parameter "TAG" is defined to represent the Git ref in repo1 ("webapp").

### **Step 3: Run Deployment Job**

1. Trigger a new deployment job in Jenkins.
2. Observe that Jenkins clones the Git SHA from repo2 ("jenkinsfiles").
3. Note that the "TAG" parameter represents the Git ref in repo1 ("webapp").

### **Step 4: Specify Jenkins Association**

1. Navigate to the Jenkins association settings in the DevLake project.
2. Specify the input parameter ("TAG") that defines the Git tag/SHA for accurate metrics.

### **Step 5: Measure DORA Metrics**

1. With the Jenkins association configured, metrics now reflect the "business-code commit" from repo1 ("webapp").
2. Measure DORA metrics accurately for deployments associated with the specified input parameter.

---

## **Conclusion**
By following these steps, you've successfully set up a Jenkins webhook for separated repositories, ensuring that DORA metrics are measured based on the correct business logic commit. This approach addresses scenarios where Jenkins jobs clone from a different repository than the one linked to GitHub, providing more accurate and meaningful metrics for your DevOps processes.

## **Additional Note**
If we could specify in the Jenkins association which input parameter defines the git tag/sha, then we would get correct metrics. This capability allows for more flexibility in associating Jenkins parameters with the appropriate repository, ensuring precise metric measurements.

---