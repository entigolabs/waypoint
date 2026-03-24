import { Card, Table, Typography } from 'antd';
import { EmsCategory } from '../client';
import React, { useState } from 'react';
import styles from './DashboardView.module.scss';

const { Title } = Typography;

const emsCategoryColumns = [
    { title: 'ID', dataIndex: 'id', key: 'id' },
    { title: 'Name', dataIndex: 'name', key: 'name' },
];

type Props = {
    emsCategories: EmsCategory[];
};

export const EmsCategoryTable: React.FC<Props> = ({ emsCategories }) => {
    const [pageSize, setPageSize] = useState(20);

    return (
        <>
            <Title level={ 4 } style={ { marginBottom: 20 } }>
                EMS Categories
            </Title>
            <Card className={ styles.card }>
                <Table
                    dataSource={ emsCategories }
                    columns={ emsCategoryColumns }
                    rowKey="id"
                    size="small"
                    pagination={ { pageSize, showSizeChanger: true, onShowSizeChange: (_, size) => setPageSize(size) } }
                    scroll={ { x: 'max-content' } }
                />
            </Card>
        </>
    );
};
