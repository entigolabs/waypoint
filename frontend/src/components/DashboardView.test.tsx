import { render, screen } from '@testing-library/react';
import { axe } from 'jest-axe';
import { expect, test, vi } from 'vitest';
import { DashboardView } from './DashboardView';

vi.mock('../client', () => ({
    getCoreCategories: vi.fn().mockResolvedValue({ data: { data: [], metadata: { total: 0 } }, error: undefined }),
    getCoreEmsCategories: vi.fn().mockResolvedValue({ data: { data: [], metadata: { total: 0 } }, error: undefined }),
    getCoreEmsThemes: vi.fn().mockResolvedValue({ data: { data: [], metadata: { total: 0 } }, error: undefined }),
}));

test('DashboardView loading state has no accessibility violations', async () => {
    const { container } = render(<DashboardView />);
    const results = await axe(container);
    expect(results).toHaveNoViolations();
});

test('DashboardView success state has no accessibility violations', async () => {
    const { container } = render(<DashboardView />);
    await screen.findByText('Categories');
    const results = await axe(container);
    expect(results).toHaveNoViolations();
});
