import { styled } from '@mui/material';
import React from 'react';

const StyledCard = styled('div')<React.HTMLAttributes<HTMLDivElement>>(({ theme }) => ({
  padding: theme.spacing(2),
  borderRadius: theme.shape.borderRadius * 2.5,
  backgroundColor: theme.palette.background.default,
})) as React.ComponentType<React.HTMLAttributes<HTMLDivElement>>;

export default StyledCard;