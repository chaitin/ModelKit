import React, { useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { ModelProvider } from '@/providers/ModelProvider';
import Model from '@/pages/model';
import ModelAdd from '@/pages/model/components/modelModal';
import { ConstsModelType } from '@/api/types';

function App() {
  const [modalOpen, setModalOpen] = useState(true); // 启动时就显示模态框

  const handleModalClose = () => {
    setModalOpen(false);
  };

  const handleRefresh = () => {
    // 这里可以添加刷新逻辑，暂时为空
  };

  return (
    <ModelProvider>
      <Router>
        <div className="App">
          <Routes>
            <Route path="/model" element={<Model />} />
            <Route path="/" element={<Navigate to="/model" replace />} />
          </Routes>
          
          {/* 在应用启动时就显示模型添加弹窗 */}
          <ModelAdd
            open={modalOpen}
            data={null}
            type={ConstsModelType.ModelTypeChat}
            onClose={handleModalClose}
            refresh={handleRefresh}
          />

        </div>
      </Router>
    </ModelProvider>
  );
}

export default App;