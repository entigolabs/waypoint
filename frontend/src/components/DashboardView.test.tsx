import { render } from '@testing-library/react';
import { axe } from 'jest-axe';
import { expect, test, vi } from 'vitest';
import { DashboardView } from './DashboardView';

vi.mock('./CategoryTable', () => ({ CategoryTable: () => <div>CategoryTable</div> }));
vi.mock('./EmsCategoryTable', () => ({ EmsCategoryTable: () => <div>EmsCategoryTable</div> }));
vi.mock('./EmsThemeTable', () => ({ EmsThemeTable: () => <div>EmsThemeTable</div> }));

test('DashboardView has no accessibility violations', async () => {
    const { container } = render(<DashboardView />);
    const results = await axe(container);
    expect(results).toHaveNoViolations();
});
