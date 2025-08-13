import { styled } from '@mui/material';

// @ts-ignore TS2742: Inferred type cannot be named without a reference
const StyledCard = styled('div')(({ theme }) => ({
  padding: theme!.spacing(2),
  borderRadius: theme!.shape.borderRadius * 2.5,
  backgroundColor: theme!.palette.background.default,
}));

export default StyledCard;
