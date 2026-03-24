import { Tag } from 'antd';
import { Category, getCoreCategories } from '../client';
import React from 'react';
import { DataTable } from './DataTable';

const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id' },
    { title: 'Name', dataIndex: 'name', key: 'name' },
    { title: 'Description', dataIndex: 'description', key: 'description' },
    {
        title: 'EMS IDs',
        dataIndex: 'emsIds',
        key: 'emsIds',
        render: (ids: string[]) =>
            ids.map((id) => <Tag key={ id }>{ id }</Tag>),
    },
];

export const CategoryTable: React.FC = () => (
    <DataTable<Category>
        title="Categories"
        columns={ columns }
        rowKey="id"
        fetchData={ getCoreCategories }
        errorMessage="Failed to load categories"
    />
);
