import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Tag, Select } from 'antd';
const { Option } = Select;

const Projects = () => {
  const [projects, setProjects] = useState([]);
  const [tags, setTags] = useState([]);
  const [selectedTags, setSelectedTags] = useState([]);
  
  useEffect(() => {
    // Fetch all tags
    axios.get('/api/tags')
      .then(response => {
        if (response.data.success) {
          setTags(response.data.tags);
        }
      })
      .catch(error => console.error('Error fetching tags:', error));
  }, []);
  
  useEffect(() => {
    // Fetch projects with tag filtering if needed
    const tagsQuery = selectedTags.length ? `tags=${selectedTags.join(',')}` : '';
    axios.get(`/api/projects?${tagsQuery}`)
      .then(response => {
        if (response.data.success) {
          setProjects(response.data.projects);
        }
      })
      .catch(error => console.error('Error fetching projects:', error));
  }, [selectedTags]);
  
  const handleTagFilterChange = (values) => {
    setSelectedTags(values);
  };
  
  // Render project tag labels
  const renderTags = (projectTags) => {
    return (
      <div className="project-tags">
        {projectTags.map(tag => (
          <Tag 
            key={tag.id} 
            color={tag.color}
            style={{ color: getContrastColor(tag.color) }}
          >
            {tag.name}
          </Tag>
        ))}
      </div>
    );
  };
  
  return (
    <div className="projects-page">
      <div className="projects-header">
        <h1>Projects</h1>
        
        <div className="tag-filter">
          <span>Filter by tags:</span>
          <Select
            mode="multiple"
            placeholder="Select tags to filter"
            onChange={handleTagFilterChange}
            style={{ minWidth: 200 }}
          >
            {tags.map(tag => (
              <Option key={tag.name} value={tag.name}>
                <Tag color={tag.color}>{tag.name}</Tag>
              </Option>
            ))}
          </Select>
        </div>
      </div>
      
      <div className="projects-list">
        {projects.map(project => (
          <div className="project-card" key={project.id}>
            <h2>{project.name}</h2>
            <p>{project.description}</p>
            
            {/* Add tag display */}
            {project.tags && project.tags.length > 0 && (
              <div className="project-tags-section">
                {renderTags(project.tags)}
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
};

// Helper function for text color contrast - same as in TagInput
function getContrastColor(hexColor) {
  // ...existing function...
}

export default Projects;
