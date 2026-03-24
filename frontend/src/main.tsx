import { App, ConfigProvider } from 'antd';
import React from 'react';
import ReactDOM from 'react-dom/client';
import AppRoot from './App';
import './assets/index.scss';
import { getAntdThemeConfig } from './theme';

ReactDOM.createRoot(document.getElementById('root')!).render(
    <React.StrictMode>
        <ConfigProvider theme={getAntdThemeConfig()}>
            <App>
                <AppRoot />
            </App>
        </ConfigProvider>
    </React.StrictMode>,
);
