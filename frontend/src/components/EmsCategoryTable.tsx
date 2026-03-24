import { EmsCategory, getCoreEmsCategories } from '../client';
import React from 'react';
import { DataTable } from './DataTable';

const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id' },
    { title: 'Name', dataIndex: 'name', key: 'name' },
];

export const EmsCategoryTable: React.FC = () => (
    <DataTable<EmsCategory>
        title="EMS Categories"
        columns={ columns }
        rowKey="id"
        fetchData={ getCoreEmsCategories }
        errorMessage="Failed to load EMS categories"
    />
);
