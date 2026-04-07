import { Alert, Card, Table, Typography } from 'antd';
import { ColumnsType } from 'antd/es/table';
import { useEffect, useRef, useState } from 'react';
import { Errors } from '../client';
import styles from './DashboardView.module.scss';

const { Title } = Typography;

const extractErrorInfo = (error: Errors | string, response: Response | undefined): ErrorState => {
    if (!response) {
        return {
            message: 'The request failed. The server response could not be read — this is likely caused by a CORS restriction on the API endpoint. Check the browser console for details.',
            code: undefined,
        };
    }
    if (typeof error === 'string') {
        return { message: error, code: response.status };
    }
    const message = error.errors?.map(e => e.message).join(', ') ?? 'An unknown error occurred';
    return { message, code: response.status };
};

type ErrorState = {
    message: string;
    code: number | undefined;
}

type DataTableProps<T extends object> = {
    title: string;
    columns: ColumnsType<T>;
    rowKey: string;
    fetchData: () => Promise<{ data?: { data: T[] }; error?: Errors | string; response?: Response }>;
    errorMessage: string;
}

export const DataTable = <T extends object>({ title, columns, rowKey, fetchData, errorMessage }: DataTableProps<T>) => {
    const [data, setData] = useState<T[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<ErrorState | null>(null);
    const [pageSize, setPageSize] = useState(20);
    const tableWrapperRef = useRef<HTMLDivElement>(null);

    useEffect(function makeScrollContainerFocusable() {
        const scrollEl = tableWrapperRef.current?.querySelector<HTMLElement>('.ant-table-content, .ant-table-body');
        if (scrollEl) {
            scrollEl.setAttribute('tabindex', '0');
            scrollEl.setAttribute('role', 'region');
            scrollEl.setAttribute('aria-label', `${ title } table`);
        }
    }, [title]);

    useEffect(function onComponentMountFetchData() {
        fetchData()
            .then(({ data, error, response }) => {
                if (error) {
                    setError(extractErrorInfo(error, response));
                    return;
                }
                if (!data || !Array.isArray(data.data)) {
                    setError({ message: 'The server returned a 200 response but the body was not valid JSON.', code: response?.status });
                    return;
                }
                setData(data.data);
            })
            .catch((err: Error) => setError({ message: err.message, code: undefined }))
            .finally(() => setLoading(false));
    }, [fetchData]);

    const alertMessage = error
        ? error.code !== undefined ? `${ errorMessage } (${ error.code })` : errorMessage
        : '';

    return (
        <>
            <Title level={ 3 } style={ { marginBottom: 20 } }>
                { title }
            </Title>
            <Card className={ styles.card }>
                { error ? (
                    <Alert type="error" title={ alertMessage } description={ error.message } showIcon />
                ) : (
                    <div ref={ tableWrapperRef }>
                        <Table
                            dataSource={ data }
                            columns={ columns }
                            rowKey={ rowKey }
                            size="small"
                            loading={ loading }
                            pagination={ { pageSize, showSizeChanger: true, onShowSizeChange: (_, size) => setPageSize(size) } }
                            scroll={ { x: 'max-content' } }
                        />
                    </div>
                ) }
            </Card>
        </>
    );
}
