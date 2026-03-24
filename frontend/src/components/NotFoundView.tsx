import { Button, Result } from 'antd';
import React from 'react';

export const NotFoundView: React.FC = () => (
    <Result
        status="404"
        title="404"
        subTitle="This page does not exist."
        extra={ <Button type="primary" href="/">Go to dashboard</Button> }
    />
);
