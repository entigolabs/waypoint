import { render, screen } from '@testing-library/react';
import { axe } from 'jest-axe';
import { expect, test, vi } from 'vitest';
import { EndpointView } from './EndpointView';

vi.mock('../client', () => ({
    getCoreCategories: vi.fn().mockResolvedValue({ data: { data: [], metadata: { total: 0 } }, error: undefined }),
    getCoreEmsCategories: vi.fn().mockResolvedValue({ data: { data: [], metadata: { total: 0 } }, error: undefined }),
    getCoreEmsThemes: vi.fn().mockResolvedValue({ data: { data: [], metadata: { total: 0 } }, error: undefined }),
}));

test('EndpointView loading state has no accessibility violations', async () => {
    const { container } = render(<EndpointView />);
    const results = await axe(container);
    expect(results).toHaveNoViolations();
});

test('EndpointView success state has no accessibility violations', async () => {
    const { container } = render(<EndpointView />);
    await screen.findByText('Categories');
    const results = await axe(container);
    expect(results).toHaveNoViolations();
});
