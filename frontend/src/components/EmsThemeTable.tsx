import { Tag } from 'antd';
import { EmsTheme, getCoreEmsThemes } from '../client';
import React from 'react';
import { DataTable } from './DataTable';

const columns = [
    { title: 'Code', dataIndex: 'code', key: 'code' },
    { title: 'Datasets Count', dataIndex: 'datasetsCount', key: 'datasetsCount' },
    {
        title: 'Translations',
        dataIndex: 'translations',
        key: 'translations',
        render: (translations: EmsTheme['translations']) =>
            translations.map((t) => (
                <Tag key={ t.language }>{ t.language }: { t.value }</Tag>
            )),
    },
    {
        title: 'EMS IDs',
        dataIndex: 'emsIds',
        key: 'emsIds',
        render: (ids: string[]) =>
            ids.map((id) => <Tag key={ id }>{ id }</Tag>),
    },
];

export const EmsThemeTable: React.FC = () => (
    <DataTable<EmsTheme>
        title="EMS Themes"
        columns={ columns }
        rowKey="code"
        fetchData={ getCoreEmsThemes }
        errorMessage="Failed to load EMS themes"
    />
);
