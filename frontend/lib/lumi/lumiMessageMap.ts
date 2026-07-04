import { lumiTriggerMessageTemplates } from './lumiMessageTemplates';
import type { LumiMessageType } from './lumiTypes';

const entries = Object.values(lumiTriggerMessageTemplates);

export const lumiMessageSamples: Record<LumiMessageType, string[]> = {
  summary: entries.filter((entry) => entry.messageType === 'summary').map((entry) => entry.message),
  hint: entries.filter((entry) => entry.messageType === 'hint').map((entry) => entry.message),
  question: entries.filter((entry) => entry.messageType === 'question').map((entry) => entry.message),
  encourage: entries.filter((entry) => entry.messageType === 'encourage').map((entry) => entry.message),
  success: entries.filter((entry) => entry.messageType === 'success').map((entry) => entry.message),
  guide: entries.filter((entry) => entry.messageType === 'guide').map((entry) => entry.message),
  discovery: entries.filter((entry) => entry.messageType === 'discovery').map((entry) => entry.message),
};
