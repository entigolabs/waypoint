import { Alert, Card, Flex, Spin, Table, Tag, Typography } from 'antd';
import { getCoreCategories, getCoreEmsCategories, getCoreEmsThemes, Category, EmsCategory, EmsTheme } from '../client';
import { client } from '../client/client.gen';
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

function extractErrorInfo(error: unknown, response: Response | undefined): { message: string; code: number | undefined } {
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
