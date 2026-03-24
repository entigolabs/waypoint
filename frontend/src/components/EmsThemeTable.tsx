import { Card, Table, Tag, Typography } from 'antd';
import { EmsTheme } from '../client';
import React, { useState } from 'react';
import styles from './DashboardView.module.scss';

const { Title } = Typography;

const emsThemeColumns = [
    { title: 'Code', dataIndex: 'code', key: 'code' },
    { title: 'Datasets Count', dataIndex: 'datasetsCount', key: 'datasetsCount' },
    {
        title: 'EMS IDs',
        dataIndex: 'emsIds',
        key: 'emsIds',
        render: (ids: string[]) =>
            ids.map((id) => <Tag key={ id }>{ id }</Tag>),
    },
    {
        title: 'Translations',
        dataIndex: 'translations',
        key: 'translations',
        render: (translations: EmsTheme['translations']) =>
            translations.map((t) => (
                <Tag key={ t.language }>{ t.language }: { t.value }</Tag>
            )),
    },
];

type Props = {
    emsThemes: EmsTheme[];
};

export const EmsThemeTable: React.FC<Props> = ({ emsThemes }) => {
    const [pageSize, setPageSize] = useState(20);

    return (
        <>
            <Title level={ 4 } style={ { marginBottom: 20 } }>
                EMS Themes
            </Title>
            <Card className={ styles.card }>
                <Table
                    dataSource={ emsThemes }
                    columns={ emsThemeColumns }
                    rowKey="code"
                    size="small"
                    pagination={ { pageSize, showSizeChanger: true, onShowSizeChange: (_, size) => setPageSize(size) } }
                    scroll={ { x: 'max-content' } }
                />
            </Card>
        </>
    );
};
