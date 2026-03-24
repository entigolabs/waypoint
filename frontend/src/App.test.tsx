import { render, screen } from '@testing-library/react';
import { axe } from 'jest-axe';
import { expect, test, vi } from 'vitest';
import App from './App';

vi.mock('./components/DashboardView', () => ({
    DashboardView: () => <div>DashboardView</div>,
}));

test('renders DashboardView on index page', () => {
    render(<App />);
    screen.getByText('DashboardView');
});

test('renders 404 on unknown page', () => {
    const original = window.location.pathname;
    Object.defineProperty(window, 'location', { value: { ...window.location, pathname: '/unknown' }, writable: true, configurable: true });
    render(<App />);
    screen.getByText('404');
    Object.defineProperty(window, 'location', { value: { ...window.location, pathname: original }, writable: true, configurable: true });
});

test('index page has no accessibility violations', async () => {
    const { container } = render(<App />);
    const results = await axe(container);
    expect(results).toHaveNoViolations();
});
