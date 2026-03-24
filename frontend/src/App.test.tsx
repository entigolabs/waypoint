import { render } from '@testing-library/react';
import { axe } from 'jest-axe';
import { expect, test, vi } from 'vitest';
import App from './App';

vi.mock('./components/DashboardView', () => ({
    DashboardView: () => <div>DashboardView</div>,
}));

test('App has no accessibility violations', async () => {
    const { container } = render(<App />);
    const results = await axe(container);
    expect(results).toHaveNoViolations();
});
