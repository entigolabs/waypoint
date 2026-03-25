import { ConfigProvider, Flex, Layout, Select, Typography, Image, Button, Drawer, Tooltip } from 'antd';
import { MenuOutlined } from '@ant-design/icons';
import React, { useState } from 'react';
import entigoLogo from './assets/entigo.svg';
import { DashboardView } from './components/DashboardView';
import { NotFoundView } from './components/NotFoundView';
import styles from './App.module.scss';

const { Header, Content } = Layout;
const { Text } = Typography;

const fontSizeOptions = [
    { label: 'Small', value: 12 },
    { label: 'Medium', value: 14 },
    { label: 'Large', value: 16 },
    { label: 'Extra Large', value: 18 },
];

const renderFontSizeOption = (option: { label?: React.ReactNode; value?: string | number | null }) => (
    <Text style={ { fontSize: option.value as number } }>{ option.label }</Text>
);

const App: React.FC = () => {
    const [fontSize, setFontSize] = useState(14);
    const [menuOpen, setMenuOpen] = useState(false);
    const apiUrl = import.meta.env.VITE_API_ENDPOINT || "";
    const isIndexPage = window.location.pathname === '/';

    const controls = (
        <>
            <Flex gap={ 8 } align="center" style={ { minWidth: 0 } }>
                <Text style={ { whiteSpace: 'nowrap' } }>API URL:</Text>
                <Tooltip
                    title={ `${ apiUrl }/api` }
                    color="var(--custom-tooltip-bg)"
                    mouseLeaveDelay={ 0.2 }
                >
                    <Text type="secondary" ellipsis>
                        { `${ apiUrl }/api` }
                    </Text>
                </Tooltip>
            </Flex>
            <Flex gap={ 8 } align="center" style={ { whiteSpace: 'nowrap', flexShrink: 0 } }>
                <Text>Font size:</Text>
                <Select
                    aria-label="Font size"
                    value={ fontSize }
                    options={ fontSizeOptions }
                    onChange={ setFontSize }
                    optionRender={ renderFontSizeOption }
                    labelRender={ renderFontSizeOption }
                    style={ { width: 150 } }
                />
            </Flex>
        </>
    );

    return (
        <ConfigProvider theme={ { token: { fontSize, fontSizeSM: fontSize } } }>
            <Layout className={ styles.layout }>
                <Header className={ styles.header }>
                    <a href="/" className={ styles.logoArea }>
                        <Image src={ entigoLogo } className={ styles.logo } alt="Entigo" preview={ false } />
                    </a>
                    <Flex gap={ 16 } align="center" className={ styles.desktopControls }>
                        { controls }
                    </Flex>
                    <Button
                        type="text"
                        icon={ <MenuOutlined /> }
                        className={ styles.hamburger }
                        onClick={ () => setMenuOpen(true) }
                        aria-label="Open menu" />
                    <Drawer
                        title="Settings"
                        placement="right"
                        open={ menuOpen }
                        onClose={ () => setMenuOpen(false) }
                        size={ 350 }
                    >
                        <Flex vertical gap={ 16 }>
                            { controls }
                        </Flex>
                    </Drawer>
                </Header>
                <Content className={ styles.content }>
                    { isIndexPage ? <DashboardView /> : <NotFoundView /> }
                </Content>
            </Layout>
        </ConfigProvider>
    );
};

export default App;
