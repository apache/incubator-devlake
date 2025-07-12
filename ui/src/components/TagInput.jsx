import React, { useState, useEffect } from 'react';
import { Tag, Input, Tooltip } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import axios from 'axios';

const TagInput = ({ value = [], onChange, maxTags = 10 }) => {
  const [tags, setTags] = useState(value);
  const [inputVisible, setInputVisible] = useState(false);
  const [inputValue, setInputValue] = useState('');
  const [allTags, setAllTags] = useState([]);
  const inputRef = React.useRef(null);

  useEffect(() => {
    // Fetch all available tags
    axios.get('/api/tags')
      .then(response => {
        if (response.data.success) {
          setAllTags(response.data.tags);
        }
      })
      .catch(error => console.error('Error fetching tags:', error));
  }, []);

  useEffect(() => {
    if (inputVisible) {
      inputRef.current?.focus();
    }
  }, [inputVisible]);

  useEffect(() => {
    setTags(value);
  }, [value]);

  const handleClose = (removedTag) => {
    const newTags = tags.filter(tag => tag.id !== removedTag.id);
    setTags(newTags);
    onChange?.(newTags);
  };

  const showInput = () => {
    setInputVisible(true);
  };

  const handleInputChange = (e) => {
    setInputValue(e.target.value);
  };

  const handleInputConfirm = () => {
    if (inputValue && !tags.some(tag => tag.name.toLowerCase() === inputValue.toLowerCase())) {
      // Check if tag already exists in system
      const existingTag = allTags.find(t => t.name.toLowerCase() === inputValue.toLowerCase());
      
      if (existingTag) {
        // Use existing tag
        const newTags = [...tags, existingTag];
        setTags(newTags);
        onChange?.(newTags);
      } else {
        // Create new tag
        axios.post('/api/tags', { 
          name: inputValue, 
          color: `#${Math.floor(Math.random()*16777215).toString(16)}` 
        })
          .then(response => {
            if (response.data.success) {
              const newTag = response.data.tag;
              const newTags = [...tags, newTag];
              setTags(newTags);
              onChange?.(newTags);
              // Add to allTags list
              setAllTags([...allTags, newTag]);
            }
          })
          .catch(error => console.error('Error creating tag:', error));
      }
    }
    
    setInputVisible(false);
    setInputValue('');
  };

  return (
    <>
      {tags.map(tag => {
        return (
          <Tag
            key={tag.id}
            closable
            style={{ color: getContrastColor(tag.color), backgroundColor: tag.color }}
            onClose={() => handleClose(tag)}
          >
            {tag.name}
          </Tag>
        );
      })}
      
      {inputVisible && (
        <Input
          ref={inputRef}
          type="text"
          size="small"
          className="tag-input"
          value={inputValue}
          onChange={handleInputChange}
          onBlur={handleInputConfirm}
          onPressEnter={handleInputConfirm}
          style={{ width: 78 }}
        />
      )}
      
      {!inputVisible && tags.length < maxTags && (
        <Tag onClick={showInput} className="site-tag-plus">
          <PlusOutlined /> New Tag
        </Tag>
      )}
    </>
  );
};

// Helper function for text color contrast
function getContrastColor(hexColor) {
  // Convert hex to RGB
  const r = parseInt(hexColor.slice(1, 3), 16);
  const g = parseInt(hexColor.slice(3, 5), 16);
  const b = parseInt(hexColor.slice(5, 7), 16);
  
  // Calculate brightness
  const brightness = (r * 299 + g * 587 + b * 114) / 1000;
  
  // Return black or white depending on brightness
  return brightness > 128 ? '#000000' : '#FFFFFF';
}

export default TagInput;
