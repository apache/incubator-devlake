import React from 'react';
import { Form, Input, Button } from 'antd';
import TagInput from '../../components/TagInput';

const ProjectForm = ({ project, onSubmit, ...props }) => {
  const [form] = Form.useForm();

  const handleSubmit = (values) => {
    onSubmit(values);
  };

  return (
    <Form
      form={form}
      initialValues={project}
      onFinish={handleSubmit}
      {...props}
    >
      <Form.Item
        name="name"
        label="Project Name"
        rules={[{ required: true, message: 'Please input the project name!' }]}
      >
        <Input />
      </Form.Item>

      <Form.Item
        name="description"
        label="Description"
        rules={[{ required: true, message: 'Please input the project description!' }]}
      >
        <Input.TextArea />
      </Form.Item>

      <Form.Item
        name="tags"
        label="Tags"
        tooltip="Add tags to organize and filter projects"
      >
        <TagInput />
      </Form.Item>

      <Form.Item>
        <Button type="primary" htmlType="submit">
          Submit
        </Button>
      </Form.Item>
    </Form>
  );
};

export default ProjectForm;
