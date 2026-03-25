import { render, screen } from '@testing-library/react';
import { axe } from 'jest-axe';
import { expect, test, vi } from 'vitest';
import { DataTable } from './DataTable';

const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id' },
    { title: 'Name', dataIndex: 'name', key: 'name' },
];

type TestRow = { id: string; name: string };

const defaultProps = {
    title: 'Test Table',
    columns,
    rowKey: 'id',
    errorMessage: 'Failed to load test data',
};

test('renders title immediately in loading state', () => {
    const fetchData = vi.fn(() => new Promise<never>(() => { }));
    render(<DataTable<TestRow> { ...defaultProps } fetchData={ fetchData } />);
    screen.getByText('Test Table');
});

test('renders rows after successful fetch', async () => {
    const fetchData = vi.fn().mockResolvedValue({
        data: { data: [{ id: '1', name: 'Item One' }] },
        error: undefined,
    });
    render(<DataTable<TestRow> { ...defaultProps } fetchData={ fetchData } />);
    await screen.findByText('Item One');
});

test('shows error message with HTTP code when fetch returns an error response', async () => {
    const fetchData = vi.fn().mockResolvedValue({
        data: undefined,
        error: 'Not Found',
        response: { status: 404 } as Response,
    });
    render(<DataTable<TestRow> { ...defaultProps } fetchData={ fetchData } />);
    await screen.findByText('Failed to load test data (404)');
    await screen.findByText('Not Found');
});

test('shows CORS hint when fetch returns error without a response object', async () => {
    const fetchData = vi.fn().mockResolvedValue({
        data: undefined,
        error: 'Network error',
        response: undefined,
    });
    render(<DataTable<TestRow> { ...defaultProps } fetchData={ fetchData } />);
    await screen.findByText('Failed to load test data');
    await screen.findByText(/CORS/);
});

test('shows error when fetch promise rejects', async () => {
    const fetchData = vi.fn().mockRejectedValue(new Error('Unexpected failure'));
    render(<DataTable<TestRow> { ...defaultProps } fetchData={ fetchData } />);
    await screen.findByText('Failed to load test data');
    await screen.findByText('Unexpected failure');
});

test('shows error messages from JSON error body', async () => {
    const fetchData = vi.fn().mockResolvedValue({
        data: undefined,
        error: { errors: [{ code: 'InternalServerError', message: 'Internal Server Error' }] },
        response: { status: 500 } as Response,
    });
    render(<DataTable<TestRow> { ...defaultProps } fetchData={ fetchData } />);
    await screen.findByText('Failed to load test data (500)');
    await screen.findByText('Internal Server Error');
});

test('loading state has no accessibility violations', async () => {
    const fetchData = vi.fn(() => new Promise<never>(() => { }));
    const { container } = render(<DataTable<TestRow> { ...defaultProps } fetchData={ fetchData } />);
    const results = await axe(container);
    expect(results).toHaveNoViolations();
});

test('success state has no accessibility violations', async () => {
    const fetchData = vi.fn().mockResolvedValue({
        data: { data: [{ id: '1', name: 'Item One' }] },
        error: undefined,
    });
    const { container } = render(<DataTable<TestRow> { ...defaultProps } fetchData={ fetchData } />);
    await screen.findByText('Item One');
    const results = await axe(container);
    expect(results).toHaveNoViolations();
});

test('error state has no accessibility violations', async () => {
    const fetchData = vi.fn().mockResolvedValue({
        data: undefined,
        error: 'Internal Server Error',
        response: { status: 500 } as Response,
    });
    const { container } = render(<DataTable<TestRow> { ...defaultProps } fetchData={ fetchData } />);
    await screen.findByText('Failed to load test data (500)');
    const results = await axe(container);
    expect(results).toHaveNoViolations();
});
