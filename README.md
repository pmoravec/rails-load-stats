# rails-load-stats

Rails-load-stats is a simple `bash` script that processes a logfile of any Ruby on Rails app to analyze where the load to the app comes from. It produces statistics with how many requests against each method were raised, how much (min/max/avg/sum) time these took, and what percentage of overall execution time was spent by processing the different types of requests. Example output:

    there were 3196 requests taking 782320 ms (i.e. 0.22 hours, i.e. 0.01 days) in summary
    
    type						count	min	max	avg	mean	sum		percentage
    --------------------------------------------------------------------------------------------------------------------
    HostsController#externalNodes                  	220	233	969	471	429	103832		13.27 %
    HostsController#facts                          	147	29	2912	671	777	98675		12.61 %
    DashboardController#show                       	84	52	11378	881	318	74057		9.47 %
    JobInvocationsController#show                  	573	52	508	88	77	50705		6.48 %
    HostsController#get_power_state                	154	24	3033	267	69	41150		5.26 %
    SyncManagementController#index                 	4	8314	10158	9367	9188	37470		4.79 %
    ProductsController#index                       	2	13196	17166	15181	13196	30362		3.88 %
    CandlepinProxiesController#get                 	432	35	583	69	52	30124		3.85 %
    ..
    ..
    ..ontroller#puppet_environment_for_content_view	1	30	30	30	30	30		0.00 %


### Usage:
./analyze.sh logfile [sort-results]

where (optional) sort-results can be:
  2: sort by count
  3: sort by min time
  4: sort by max time
  5: sort by avg time
  6: sort by mean time
  7: sort by sum time / percentage

### Requirements:
`bash` version >= 4
