import { styled } from '@mui/material';
import type { Theme } from '@mui/material/styles';
import type { ComponentType } from 'react';

const StyledCard: ComponentType<any> = styled('div')<{}>(
  ({ theme }: { theme: Theme }) => ({
    padding: theme.spacing(2),
    borderRadius: (theme.shape.borderRadius as number) * 2.5,
    backgroundColor: theme.palette.background.default,
  })
);

export default StyledCard;
