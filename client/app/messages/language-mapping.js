import messagesEN from 'en.json'
import messageZH_HANS from 'zh_HANS.json'
import messageDE from 'de.json'
import messageES from 'es.json'
import messageJA from 'ja.json'
import messageFR from 'fr.json'
import enLocaleData from 'react-intl/locale-data/en';
import zhLocaleData from 'react-intl/locale-data/zh';
import deLocaleData from 'react-intl/locale-data/de';
import esLocaleData from 'react-intl/locale-data/es';
import jaLocaleData from 'react-intl/locale-data/ja';
import frLocaleData from 'react-intl/locale-data/fr';

export const languageMapping = {
  en:        {
    messages:   messagesEN,
    name:       'English',
    locale:     'en',
    localeData: enLocaleData,
  },
  'zh-Hans': {
    messages:   messageZH_HANS,
    name:       'Chinese (simplified)',
    locale:     'zh-Hans',
    localeData: zhLocaleData,
  },
  es:        {
    messages:   messageES,
    name:       'Spanish',
    locale:     'es',
    localeData: esLocaleData,
  },
  de:        {
    messages:   messageDE,
    name:       'German',
    locale:     'de',
    localeData: deLocaleData,
  },
  ja:        {
    messages:   messageJA,
    name:       'Japanese',
    locale:     'ja',
    localeData: jaLocaleData,
  },
  fr:        {
    messages:   messageFR,
    name:       'French',
    locale:     'fr',
    localeData: frLocaleData,
  },
};

