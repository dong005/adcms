import { requestClient } from '#/api/request';

export interface TableInfo {
  name: string;
  engine: string;
  rows: number;
  data_size: string;
  index_size: string;
  comment: string;
}

export interface ColumnInfo {
  name: string;
  type: string;
  nullable: string;
  key: string;
  default: string | null;
  comment: string;
}

export function getTableList() {
  return requestClient.get<TableInfo[]>('/database/tables');
}

export function getTableColumns(table: string) {
  return requestClient.get<ColumnInfo[]>(`/database/tables/${table}/columns`);
}
