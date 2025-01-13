import { type SavingStates } from './common';

export interface Pages {
  title?: string;
  slug?: string;
  content?: string;
  visible?: boolean;
  contentType?: string;
  savingState?: SavingStates;
  onSubmit?: (
    data: any,
    done: (savingState: SavingStates) => void,
    error: (error: any) => void
  ) => void;
}
