import { ThemeConfig, theme } from 'antd';

const rgba = (rgb: string, alpha = '1') => `rgba(${ rgb }, ${ alpha })`;

export const c = {
    colorBgBase: '#141414' as const,
    siderBg: '#16181C' as const,
    darkSelectedBlueRgb: '44, 125, 253' as const,
    darkItemColor: '#ADAEB0' as const,
    bgSpotlight: '36, 36, 36' as const,
    borderAxis: '#2c2d30' as const,
    whiteColor: '#FFFFFF' as const,
    whiteColor25: '#505153' as const,
    whiteColor85: '#DCDCDC' as const,
    electricPurple: '#BB64FF' as const,
};

const cc = {
    darkItemSelectedColor: rgba(c.darkSelectedBlueRgb),
    darkItemSelectedBorderColor: rgba(c.darkSelectedBlueRgb, '.4'),
    darkItemSelectedBg: rgba(c.darkSelectedBlueRgb, '.15'),
};

export const getAntdThemeConfig = (): ThemeConfig => ({
    algorithm: theme.darkAlgorithm,
    cssVar: {},
    token: {
        colorBgBase: c.colorBgBase,
        colorBgLayout: c.colorBgBase,
        colorBgSpotlight: rgba(c.bgSpotlight),
        colorText: c.whiteColor85,
        colorTextHeading: c.whiteColor85,
        colorBorderSecondary: '#1C293D',
        colorBgContainer: 'transparent',
        colorPrimary: '#1677FF',
        colorPrimaryHover: '#1677FF',
        colorBgElevated: rgba(c.bgSpotlight),
        colorSplit: '#FFFFFF0F',
        fontFamily: "-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif",
    },
    components: {
        Layout: {
            headerHeight: 56,
            headerBg: c.siderBg,
            siderBg: c.siderBg,
            bodyBg: c.colorBgBase,
        },
        Menu: {
            itemBg: c.siderBg,
            darkItemSelectedBg: cc.darkItemSelectedBg,
            darkItemSelectedColor: cc.darkItemSelectedColor,
            itemHoverBg: cc.darkItemSelectedBg,
            darkItemColor: c.darkItemColor,
            darkItemHoverColor: c.darkItemColor,
            itemColor: c.darkItemColor,
            itemHoverColor: c.darkItemColor,
        },
        Table: {
            borderColor: '#24262A',
        },
        Card: {
            colorBgContainer: '#1a1a1a',
        },
    },
});
