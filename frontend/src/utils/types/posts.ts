import { type SavingStates } from "./common";

export interface Posts {
  title?: string;
  slug?: string;
  tags?: string[];
  excerpt?: string;
  content?: string;
  visible?: boolean;
  publish?: string;
  publishedAt?: string;
  savingState?: SavingStates;
  onSubmit?: (data: any, done: (savingState: SavingStates) => void, error: (error: any) => void) => void;
}
