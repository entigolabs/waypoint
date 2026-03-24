import { ConfigProvider, Flex, Layout, Select, Typography } from 'antd';
import React, { useState } from 'react';
import entigoLogo from './assets/entigo.svg';
import { EndpointView } from './components/EndpointView';
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

    return (
        <ConfigProvider theme={ { token: { fontSize, fontSizeSM: fontSize } } }>
            <Layout className={ styles.layout }>
                <Header className={ styles.header }>
                    <div className={ styles.logoArea }>
                        <img src={ entigoLogo } className={ styles.logo } alt="Entigo" />
                    </div>
                    <Flex gap={ 16 } align="center">
                        <Text type="secondary">
                            { import.meta.env.VITE_API_ENDPOINT ?? '' }
                        </Text>
                        <Flex gap={ 8 } align="center">
                            <Text>Font size:</Text>
                            <Select
                                aria-label="Font size"
                                value={ fontSize }
                                options={ fontSizeOptions }
                                onChange={ setFontSize }
                                optionRender={ renderFontSizeOption }
                                labelRender={ renderFontSizeOption }
                                style={ { width: 130 } } />
                        </Flex>
                    </Flex>
                </Header>
                <Content className={ styles.content }>
                    <EndpointView />
                </Content>
            </Layout>
        </ConfigProvider>
    );
};

export default App;
