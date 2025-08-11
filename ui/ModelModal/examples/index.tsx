import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import { AppBar, Toolbar, Typography, Container, Box, Button } from '@mui/material';
import BasicUsageExample from './basic-usage';
import AdvancedConfigExample from './advanced-config';

export const ExamplesApp: React.FC = () => {
  return (
    <Router>
      <AppBar position="static">
        <Toolbar>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            ModelModal 示例应用
          </Typography>
          <Box sx={{ display: 'flex', gap: 1 }}>
            <Button color="inherit" component={Link} to="/">
              基本使用
            </Button>
            <Button color="inherit" component={Link} to="/advanced">
              高级配置
            </Button>
          </Box>
        </Toolbar>
      </AppBar>

      <Container>
        <Routes>
          <Route path="/" element={<BasicUsageExample />} />
          <Route path="/advanced" element={<AdvancedConfigExample />} />
        </Routes>
      </Container>
    </Router>
  );
};

export default ExamplesApp; 