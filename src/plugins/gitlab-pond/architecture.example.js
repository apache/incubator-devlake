const config = {
  groups_path: {
    url: 'https://gitlab.com/api/v3/proj/:groupId/mr/:projectId/notes',
    collection: 'merge_requests'
  },
  mr_path: {
    url: 'https://gitlab.com/api/v3/proj/:projectId/mr/:mrId/notes',
    collection: 'merge_requests'
  },
  proj_path: {
    url: 'https://gitlab.com/api/v3/proj/:projectId',
    collection: 'projects'
  }
}


let buildApiUrlsForProjects = (projIdsIWant)=>{
  let apiUrls = []
  projIdsIWant.forEach(projId => {
    
    // format url
    apiUrl = config.proj_path.replace(':projectId', projId)

    // collection is needed to know which mongo table to store in
    apiUrls.push({url: apiUrl, collection: 'projects'})
  })

  return apiUrls
}

let buildApiUrlsForMRs = (projIdsIWant)=>{
  let apiUrls = []
  projIdsIWant.forEach(projId => {
    //TODO: get all MR ids for project with a gitlab API call?
    mrIds = [1,2,3,4,5]
    
    // loop through app MR ids to create the urls
    mrIds.forEach(mrId => {
      // call gitlab api for path url/proj/:projectId/mr/:mrId/notes
      apiUrl = config.mr_path.replace(':projectId', projId).replace(':mrId', mrId)
  
      // collection is needed to know which mongo table to store in
      apiUrls.push({url: apiUrl, collection: 'merge_requests'})
    })
  
  })

  return apiUrls
}


let gatherAllData = (urls) => {
  urls.forEach(url => {
    // axios post to url (url.url)
    // determine collection (url.collection)
    // save response to mongo db
  })
}



//get all project ids that the user cares about from a config? Or maybe all the pojects they have access to?
projIdsIWant = [12, 13, 55]


let allUrls = []
allUrls.push(buildApiUrlsForMRs(projIdsIWant))
allUrls.push(buildApiUrlsForProjects(projIdsIWant))

gatherAllData(allUrls.flat())

console.log('JON >>> allUrls', allUrls.flat())