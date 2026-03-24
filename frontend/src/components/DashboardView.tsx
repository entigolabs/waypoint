import { Alert, Flex, Spin } from 'antd';
import { getCoreCategories, getCoreEmsCategories, getCoreEmsThemes, Category, EmsCategory, EmsTheme } from '../client';
import { client } from '../client/client.gen';
import React, { useEffect, useState } from 'react';
import styles from './DashboardView.module.scss';
import { CategoryTable } from './CategoryTable';
import { EmsCategoryTable } from './EmsCategoryTable';
import { EmsThemeTable } from './EmsThemeTable';

type FetchStatus = 'loading' | 'success' | 'error';

type State = {
    categories: Category[];
    emsCategories: EmsCategory[];
    emsThemes: EmsTheme[];
    status: FetchStatus;
    error: string;
    errorCode: number | undefined;
}

const extractErrorInfo = (error: unknown, response: Response | undefined): { message: string; code: number | undefined } => {
    if (!response) {
        return {
            message: 'The request failed. The server response could not be read — this is likely caused by a CORS restriction on the API endpoint. Check the browser console for details.',
            code: undefined,
        };
    }
    const code = response.status;
    const message = typeof error === 'string' ? error : error instanceof Error ? error.message : String(error);
    return { message, code };
}

client.setConfig({ baseUrl: import.meta.env.VITE_API_ENDPOINT || '' });

export const DashboardView: React.FC = () => {
    const [state, setState] = useState<State>({
        categories: [],
        emsCategories: [],
        emsThemes: [],
        status: 'loading',
        error: '',
        errorCode: undefined,
    });

    useEffect(() => {
        Promise.all([
            getCoreCategories(),
            getCoreEmsCategories(),
            getCoreEmsThemes(),
        ])
            .then(([cats, emsCats, emsThemes]) => {
                const err = cats.error ?? emsCats.error ?? emsThemes.error;
                if (err) {
                    const response = (cats.response ?? emsCats.response ?? emsThemes.response) as Response | undefined;
                    const { message, code } = extractErrorInfo(err, response);
                    setState((prev) => ({
                        ...prev,
                        status: 'error',
                        error: message,
                        errorCode: code,
                    }));
                    return;
                }
                setState({
                    categories: cats.data!.data,
                    emsCategories: emsCats.data!.data,
                    emsThemes: emsThemes.data!.data,
                    status: 'success',
                    error: '',
                    errorCode: undefined,
                });
            })
            .catch((err: unknown) => {
                setState((prev) => ({
                    ...prev,
                    status: 'error',
                    error: String(err),
                    errorCode: undefined,
                }));
            });
    }, []);

    if (state.status === 'loading') {
        return (
            <Flex justify="center" style={ { padding: 48 } }>
                <Spin size="large" />
            </Flex>
        );
    }

    if (state.status === 'error') {
        const message = state.errorCode !== undefined
            ? `Failed to load data (${ state.errorCode })`
            : 'Failed to load data';

        return (
            <div className={ styles.wrapper }>
                <Alert type="error" message={ message } description={ state.error } showIcon />
            </div>
        );
    }

    return (
        <div className={ styles.wrapper }>
            <CategoryTable categories={ state.categories } />
            <EmsCategoryTable emsCategories={ state.emsCategories } />
            <EmsThemeTable emsThemes={ state.emsThemes } />
        </div>
    );
};
