import { client } from '../client/client.gen';
import React from 'react';
import styles from './DashboardView.module.scss';
import { CategoryTable } from './CategoryTable';
import { EmsCategoryTable } from './EmsCategoryTable';
import { EmsThemeTable } from './EmsThemeTable';
import { Typography } from 'antd';

client.setConfig({ baseUrl: import.meta.env.VITE_API_ENDPOINT ? `${ import.meta.env.VITE_API_ENDPOINT }/api` : '/api' });

export const DashboardView: React.FC = () => {
    return (
        <div className={ styles.wrapper }>
            <Typography.Title level={ 1 }>Dashboard</Typography.Title>
            <CategoryTable />
            <EmsCategoryTable />
            <EmsThemeTable />
        </div>
    );
};
