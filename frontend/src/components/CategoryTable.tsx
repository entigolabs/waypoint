import { Card, Table, Tag, Typography } from 'antd';
import { Category } from '../client';
import React, { useState } from 'react';
import styles from './DashboardView.module.scss';

const { Title } = Typography;

const categoryColumns = [
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

type Props = {
    categories: Category[];
};

export const CategoryTable: React.FC<Props> = ({ categories }) => {
    const [pageSize, setPageSize] = useState(20);

    return (
        <>
            <Title level={ 4 } style={ { marginBottom: 20 } }>
                Categories
            </Title>
            <Card className={ styles.card }>
                <Table
                    dataSource={ categories }
                    columns={ categoryColumns }
                    rowKey="id"
                    size="small"
                    pagination={ { pageSize, showSizeChanger: true, onShowSizeChange: (_, size) => setPageSize(size) } }
                    scroll={ { x: 'max-content' } }
                />
            </Card>
        </>
    );
};
