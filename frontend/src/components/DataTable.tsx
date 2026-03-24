import { Alert, Card, Table, Typography } from 'antd';
import { ColumnsType } from 'antd/es/table';
import { useEffect, useState } from 'react';
import styles from './DashboardView.module.scss';

const { Title } = Typography;

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
};

type ErrorState = {
    message: string;
    code: number | undefined;
}

type DataTableProps<T extends object> = {
    title: string;
    columns: ColumnsType<T>;
    rowKey: string;
    fetchData: () => Promise<{ data?: { data: T[] }; error?: unknown; response?: Response }>;
    errorMessage: string;
}

export function DataTable<T extends object>({ title, columns, rowKey, fetchData, errorMessage }: DataTableProps<T>) {
    const [data, setData] = useState<T[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<ErrorState | null>(null);
    const [pageSize, setPageSize] = useState(20);

    useEffect(function onComponentMountFetchData() {
        fetchData()
            .then(({ data, error, response }) => {
                if (error) {
                    setError(extractErrorInfo(error, response));
                    return;
                }
                setData(data!.data);
            })
            .catch((err: unknown) => setError({ message: String(err), code: undefined }))
            .finally(() => setLoading(false));
    }, [fetchData]);

    const alertMessage = error
        ? error.code !== undefined ? `${ errorMessage } (${ error.code })` : errorMessage
        : '';

    return (
        <>
            <Title level={ 4 } style={ { marginBottom: 20 } }>
                { title }
            </Title>
            <Card className={ styles.card }>
                { error ? (
                    <Alert type="error" message={ alertMessage } description={ error.message } showIcon />
                ) : (
                    <Table
                        dataSource={ data }
                        columns={ columns }
                        rowKey={ rowKey }
                        size="small"
                        loading={ loading }
                        pagination={ { pageSize, showSizeChanger: true, onShowSizeChange: (_, size) => setPageSize(size) } }
                        scroll={ { x: 'max-content' } }
                    />
                ) }
            </Card>
        </>
    );
}
