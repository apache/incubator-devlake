## How to create recurring pipelines?

Once you've verified a pipeline works well, mostly likely you'll want to run that pipeline periodically to keep data fresh, and DevLake's pipeline blueprint feature have got you covered.


1. Click 'Create Pipeline Run' and 
  - Toggle the plugins you'd like to run, here we use GitHub and GitExtractor plugin as an example
  - Toggle on Automate Pipeline
    ![image](https://user-images.githubusercontent.com/14050754/163596590-484e4300-b17e-4119-9818-52463c10b889.png)


2. Click 'Add Blueprint'. Fill in the form and 'Save Blueprint'.
    
    - **NOTE**: That the schedule syntax is standard unix cron syntax, [Crontab.guru](https://crontab.guru/) could be a useful reference
    - **IMPORANT**: The scheduler is running under `UTC` timezone. If you prefer data collecting happens at 3am NewYork(UTC-04:00) every day, use **Custom Shedule** and set it to `0 7 * * *`
    
    ![image](https://user-images.githubusercontent.com/14050754/163596655-db59e154-405f-4739-89f2-7dceab7341fe.png)
    
3. Click 'Save Blueprint'.
    
4. Click 'Pipeline Blueprints', you can view and edit the new blueprint in the blueprint list.
    
    ![image](https://user-images.githubusercontent.com/14050754/163596773-4fb4237e-e3f2-4aef-993f-8a1499ca30e2.png)