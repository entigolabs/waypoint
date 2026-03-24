import { client } from '../client/client.gen';
import React from 'react';
import styles from './DashboardView.module.scss';
import { CategoryTable } from './CategoryTable';
import { EmsCategoryTable } from './EmsCategoryTable';
import { EmsThemeTable } from './EmsThemeTable';

client.setConfig({ baseUrl: import.meta.env.VITE_API_ENDPOINT || '' });

export const DashboardView: React.FC = () => {
    return (
        <div className={ styles.wrapper }>
            <CategoryTable />
            <EmsCategoryTable />
            <EmsThemeTable />
        </div>
    );
};
