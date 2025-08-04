import React, { useState } from 'react';
import Card from '@/components/card';
import { useRequest } from 'ahooks';
import { getListModel } from '@/api/Model';
import { DomainModel, ConstsModelType } from '@/api/types';
import { Stack, Box, Button, Grid2 as Grid, ButtonBase } from '@mui/material';
import StyledLabel from '@/components/label';
import { Icon, Modal, message } from '@c-x/ui';
import { addCommasToNumber } from '@/utils';
import NoData from '@/assets/images/nodata.png';
import { ModelProvider } from '../constant';
import ModelModal from './modelModal';

const ModelItem = ({
  data,
  onEdit,
  refresh,
}: {
  data: DomainModel;
  onEdit: (data: DomainModel) => void;
  refresh: () => void;
}) => {
  const onInactiveModel = () => {
    Modal.confirm({
      title: '停用模型',
      content: (
        <>
          确定要停用{' '}
          <Box component='span' sx={{ fontWeight: 700, color: 'text.primary' }}>
            {data.id}
          </Box>{' '}
          模型吗？
        </>
      ),
      okText: '停用',
      okButtonProps: {
        color: 'error',
      },
      onOk: () => {
        message.info('模型停用功能待开发');
        // putUpdateModel({
        //   id: data.id,
        //   status: ConstsModelStatus.ModelStatusInactive,
        //   provider: data.owned_by!,
        // }).then(() => {
        //   message.success('停用成功');
        //   refresh();
        // });
      },
    });
  };

  const onRemoveModel = () => {
    Modal.confirm({
      title: '删除模型',
      content: (
        <>
          确定要删除{' '}
          <Box component='span' sx={{ fontWeight: 700, color: 'text.primary' }}>
            {data.id}
          </Box>{' '}
          模型吗？
        </>
      ),
      okText: '删除',
      okButtonProps: {
        color: 'error',
      },
      onOk: () => {
        message.info('模型删除功能待开发');
        // deleteDeleteModel({ id: data.id! }).then(() => {
        //   message.success('删除成功');
        //   refresh();
        // });
      },
    });
  };

  const onActiveModel = () => {
    Modal.confirm({
      title: '激活模型',
      content: (
        <>
          确定要激活{' '}
          <Box component='span' sx={{ fontWeight: 700, color: 'text.primary' }}>
            {data.id}
          </Box>{' '}
          模型吗？
        </>
      ),
      onOk: () => {
        message.info('模型激活功能待开发');
        // putUpdateModel({
        //   id: data.id,
        //   status: ConstsModelStatus.ModelStatusActive,
        //   provider: data.owned_by!,
        // }).then(() => {
        //   message.success('激活成功');
        //   refresh();
        // });
      },
    });
  };

  return (
    <Card
      sx={{
        overflow: 'hidden',
        position: 'relative',
        transition: 'all 0.3s ease',
        borderStyle: 'solid',
        borderWidth: '1px',
        borderColor: 'success.main',
        boxShadow:
          '0px 0px 10px 0px rgba(68, 80, 91, 0.1), 0px 0px 2px 0px rgba(68, 80, 91, 0.1)',
        '&:hover': {
          boxShadow:
            'rgba(54, 59, 76, 0.3) 0px 10px 30px 0px, rgba(54, 59, 76, 0.03) 0px 0px 1px 1px',
        },
      }}
    >
      <Stack
        direction='row'
        alignItems='center'
        justifyContent='space-between'
        sx={{ height: 28 }}
      >
        <Stack direction='row' alignItems='center' gap={1}>
          <Icon
            type={
              ModelProvider[data.owned_by as keyof typeof ModelProvider]?.icon
            }
            sx={{ fontSize: 24 }}
          />
          <Stack
            direction='row'
            alignItems='center'
            gap={1}
            sx={{ fontSize: 14, minWidth: 0 }}
          >
            <Box
              sx={{
                fontWeight: 700,
                overflow: 'hidden',
                textOverflow: 'ellipsis',
                whiteSpace: 'nowrap',
              }}
            >
              {data.id || '未命名'}
            </Box>
            <Box
              sx={{
                color: 'text.tertiary',
                overflow: 'hidden',
                textOverflow: 'ellipsis',
                whiteSpace: 'nowrap',
              }}
            >
              / {data.id}
            </Box>
          </Stack>
        </Stack>
      </Stack>
    </Card>
  );
};

interface IModelCardProps {
  title: string;
  modelType: ConstsModelType;
  data: DomainModel[];
  refreshModel: () => void;
}

const ModelCard: React.FC<IModelCardProps> = ({
  title,
  modelType,
  data,
  refreshModel,
}) => {
  const [open, setOpen] = useState(false);
  const [editData, setEditData] = useState<DomainModel | null>(null);

  const onEdit = (data: DomainModel) => {
    setOpen(true);
    setEditData(data);
  };

  return (
    <Card>
      <Stack direction='row' justifyContent='space-between' alignItems='center'>
        <Box sx={{ fontWeight: 700 }}>{title}</Box>
        <Button
          variant='contained'
          color='primary'
          onClick={() => setOpen(true)}
        >
          添加模型
        </Button>
      </Stack>
      {data?.length > 0 ? (
        <Grid container spacing={2} sx={{ mt: 2 }}>
          {data.map((item) => (
            <Grid size={{ xs: 12, sm: 12, md: 12, lg: 6, xl: 4 }} key={item.id}>
              <ModelItem data={item} onEdit={onEdit} refresh={refreshModel} />
            </Grid>
          ))}
        </Grid>
      ) : (
        <Stack alignItems={'center'} sx={{ my: 2 }}>
          <img src={NoData} width={150} alt='empty' />
          <Box sx={{ color: 'error.main', fontSize: 12 }}>
            暂无模型，请先添加模型
          </Box>
        </Stack>
      )}

      <ModelModal
        open={open}
        onClose={() => {
          setOpen(false);
          setEditData(null);
        }}
        refresh={refreshModel}
        data={editData}
        type={modelType}
      />
    </Card>
  );
};

export default ModelCard;
