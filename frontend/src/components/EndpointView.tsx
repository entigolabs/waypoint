import { Alert, Card, Flex, Spin, Table, Tag, Typography } from 'antd';
import { Configuration, DefaultApi } from '@entigolabs/waypoint-api';
import type { Category, EmsCategory, EmsTheme } from '@entigolabs/waypoint-api';
import React, { useEffect, useState } from 'react';
import styles from './EndpointView.module.scss';

const { Title } = Typography;

type FetchStatus = 'loading' | 'success' | 'error';

interface State {
    categories: Category[];
    emsCategories: EmsCategory[];
    emsThemes: EmsTheme[];
    status: FetchStatus;
    error: string;
    errorCode: number | undefined;
}

async function extractErrorInfo(err: unknown): Promise<{ message: string; code: number | undefined }> {
    if (err != null && typeof err === 'object' && (err as { name?: unknown }).name === 'FetchError') {
        return {
            message: 'The request failed. The server response could not be read — this is likely caused by a CORS restriction on the API endpoint. Check the browser console for details.',
            code: undefined,
        };
    }

    if (err instanceof SyntaxError) {
        return {
            message: 'The API returned an unexpected response (not JSON). The API endpoint may not be configured correctly.',
            code: undefined,
        };
    }

    let code: number | undefined;
    let message = err instanceof Error ? err.message : String(err);

    if (err != null && typeof err === 'object') {
        const response = (err as Record<string, unknown>).response;
        if (
            response != null &&
            typeof response === 'object' &&
            'status' in response &&
            typeof (response as { status: unknown }).status === 'number'
        ) {
            code = (response as { status: number }).status;
            message = await (response as Response).text();
        }
    }

    return { message, code };
}

const api = new DefaultApi(
    new Configuration({ basePath: import.meta.env.VITE_API_ENDPOINT || undefined }),
);

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

const emsCategoryColumns = [
    { title: 'ID', dataIndex: 'id', key: 'id' },
    { title: 'Name', dataIndex: 'name', key: 'name' },
];

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

export const EndpointView: React.FC = () => {
    const [state, setState] = useState<State>({
        categories: [],
        emsCategories: [],
        emsThemes: [],
        status: 'loading',
        error: '',
        errorCode: undefined,
    });
    const [categoriesPageSize, setCategoriesPageSize] = useState(20);
    const [emsCategoriesPageSize, setEmsCategoriesPageSize] = useState(20);
    const [emsThemesPageSize, setEmsThemesPageSize] = useState(20);

    useEffect(() => {
        Promise.all([
            api.getCoreCategories(),
            api.getCoreEmsCategories(),
            api.getCoreEmsThemes(),
        ])
            .then(([cats, emsCats, emsThemes]) => {
                setState({
                    categories: cats.data,
                    emsCategories: emsCats.data,
                    emsThemes: emsThemes.data,
                    status: 'success',
                    error: '',
                    errorCode: undefined,
                });
            })
            .catch(async (err: unknown) => {
                const { message, code } = await extractErrorInfo(err);
                setState((prev) => ({
                    ...prev,
                    status: 'error',
                    error: message,
                    errorCode: code,
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
            <Title level={ 4 } style={ { marginBottom: 20 } }>
                Categories
            </Title>
            <Card className={ styles.card }>
                <Table
                    dataSource={ state.categories }
                    columns={ categoryColumns }
                    rowKey="id"
                    size="small"
                    pagination={ { pageSize: categoriesPageSize, showSizeChanger: true, onShowSizeChange: (_, size) => setCategoriesPageSize(size) } }
                    scroll={ { x: 'max-content' } }
                />
            </Card>

            <Title level={ 4 } style={ { marginBottom: 20 } }>
                EMS Categories
            </Title>
            <Card className={ styles.card }>
                <Table
                    dataSource={ state.emsCategories }
                    columns={ emsCategoryColumns }
                    rowKey="id"
                    size="small"
                    pagination={ { pageSize: emsCategoriesPageSize, showSizeChanger: true, onShowSizeChange: (_, size) => setEmsCategoriesPageSize(size) } }
                    scroll={ { x: 'max-content' } }
                />
            </Card>

            <Title level={ 4 } style={ { marginBottom: 20 } }>
                EMS Themes
            </Title>
            <Card className={ styles.card }>
                <Table
                    dataSource={ state.emsThemes }
                    columns={ emsThemeColumns }
                    rowKey="code"
                    size="small"
                    pagination={ { pageSize: emsThemesPageSize, showSizeChanger: true, onShowSizeChange: (_, size) => setEmsThemesPageSize(size) } }
                    scroll={ { x: 'max-content' } }
                />
            </Card>
        </div>
    );
};
