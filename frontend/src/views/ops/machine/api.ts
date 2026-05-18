import Api, { UploadOptions } from '@/common/Api';

export const machineApi = {
    // 获取权限列表
    list: Api.newGet('/machines'),
    getByCodes: Api.newGet('/machines/simple'),
    tagList: Api.newGet('/machines/tags'),
    getMachinePwd: Api.newGet('/machines/{id}/pwd'),
    info: Api.newGet('/machines/{id}/sysinfo'),
    stats: Api.newGet('/machines/{id}/stats'),
    process: Api.newGet('/machines/{id}/process'),
    // 终止进程
    killProcess: Api.newDelete('/machines/{id}/process'),
    users: Api.newGet('/machines/{id}/users'),
    groups: Api.newGet('/machines/{id}/groups'),
    testConn: Api.newPost('/machines/test-conn'),
    // 保存按钮
    saveMachine: Api.newPost('/machines'),
    // 调整状态
    changeStatus: Api.newPut('/machines/{id}/{status}'),
    // 删除机器
    del: Api.newDelete('/machines/{id}'),
    scripts: Api.newGet('/machines/{machineId}/scripts'),
    scriptCategorys: Api.newGet('/machines/scripts/categorys'),
    runScript: Api.newGet('/machines/scripts/{scriptId}/{ac}/run'),
    saveScript: Api.newPost('/machines/{machineId}/scripts'),
    deleteScript: Api.newDelete('/machines/{machineId}/scripts/{scriptId}'),
    // 获取配置文件列表
    files: Api.newGet('/machines/{id}/files'),
    lsFile: Api.newGet('/machines/{machineId}/files/{fileId}/read-dir'),
    dirSize: Api.newGet('/machines/{machineId}/files/{fileId}/dir-size'),
    fileStat: Api.newGet('/machines/{machineId}/files/{fileId}/file-stat'),
    rmFile: Api.newPost('/machines/{machineId}/files/{fileId}/remove'),
    cpFile: Api.newPost('/machines/{machineId}/files/{fileId}/cp'),
    renameFile: Api.newPost('/machines/{machineId}/files/{fileId}/rename'),
    mvFile: Api.newPost('/machines/{machineId}/files/{fileId}/mv'),
    uploadFile: Api.newUpload('/machines/{machineId}/files/{fileId}/upload'),
    uploadFolder: Api.newPost('/machines/{machineId}/files/{fileId}/upload-folder'),
    fileContent: Api.newGet('/machines/{machineId}/files/{fileId}/read'),
    downloadFile: Api.newGet('/machines/{machineId}/files/{fileId}/download'),
    createFile: Api.newPost('/machines/{machineId}/files/{id}/create-file'),
    // 修改文件内容
    updateFileContent: Api.newPost('/machines/{machineId}/files/{id}/write'),
    // 添加文件or目录
    addConf: Api.newPost('/machines/{machineId}/files'),
    // 删除配置的文件or目录
    delConf: Api.newDelete('/machines/{machineId}/files/{id}'),
    // 机器终端操作记录列表
    termOpRecs: Api.newGet('/machines/{machineId}/term-recs'),
};

export const cronJobApi = {
    list: Api.newGet('/machine-cronjobs'),
    relateMachineIds: Api.newGet('/machine-cronjobs/machine-ids'),
    relateCronJobIds: Api.newGet('/machine-cronjobs/cronjob-ids'),
    save: Api.newPost('/machine-cronjobs'),
    delete: Api.newDelete('/machine-cronjobs/{id}'),
    run: Api.newPost('/machine-cronjobs/run/{key}'),
    execList: Api.newGet('/machine-cronjobs/execs'),
};

export const cmdConfApi = {
    list: Api.newGet('/machine/security/cmd-confs'),
    save: Api.newPost('/machine/security/cmd-confs'),
    delete: Api.newDelete('/machine/security/cmd-confs/{id}'),
};

/**
 * 获取终端 WebSocket URL
 */
export function getMachineTerminalSocketUrl(authCertName: any) {
    return `/machines/terminal/${authCertName}`;
}

/**
 * 获取 RDP WebSocket URL
 */
export function getMachineRdpSocketUrl(authCertName: any) {
    return `/api/machines/rdp/${authCertName}`;
}

/**
 * 文件上传参数
 */
export interface UploadParams {
    /** 上传ID（可选，不传则内部自动生成） */
    uploadId?: string;
    /** 机器ID */
    machineId: number;
    /** 认证证书名称 */
    authCertName: string;
    /** 协议类型 */
    protocol: number;
    /** 文件ID */
    fileId: number;
    /** 目标路径 */
    path: string;
    /** 文件名 */
    filename: string;
}

/**
 * 上传单个文件
 * @param file 文件对象
 * @param params 上传参数
 * @param options 上传选项
 * @returns { uploadId: string; abort: () => void } 返回包含 uploadId 和中止方法的对象
 */
export function uploadFile(file: File, params: UploadParams, options: UploadOptions = {}): { uploadId: string; abort: () => void } {
    // 业务层生成 uploadId
    const uploadId = params.uploadId || `upload_${Date.now()}_${Math.random().toString(36).substring(2, 11)}`;

    const formData = new FormData();
    formData.append('file', file);
    formData.append('uploadId', uploadId);
    formData.append('machineId', String(params.machineId));
    formData.append('authCertName', params.authCertName);
    formData.append('protocol', String(params.protocol));
    formData.append('fileId', String(params.fileId));
    formData.append('path', params.path);

    const { abort } = machineApi.uploadFile.upload(formData, options);

    return { uploadId, abort };
}

/**
 * 文件夹上传参数
 */
export interface FolderUploadParams {
    /** 上传ID（可选，不传则内部自动生成） */
    uploadId?: string;
    /** 机器ID */
    machineId: number;
    /** 认证证书名称 */
    authCertName: string;
    /** 协议类型 */
    protocol: number;
    /** 文件ID */
    fileId: number;
    /** 目标路径 */
    path: string;
}

/**
 * 上传文件夹(使用 /upload-folder 接口)
 * @param files 文件列表
 * @param params 上传参数
 * @param options 上传选项
 * @returns { uploadId: string; abort: () => void } 返回包含 uploadId 和中止方法的对象
 */
export function uploadFolder(files: FileList | File[], params: FolderUploadParams, options: UploadOptions = {}): { uploadId: string; abort: () => void } {
    // 业务层生成 uploadId
    const uploadId = params.uploadId || `upload_${Date.now()}_${Math.random().toString(36).substring(2, 11)}`;

    const formData = new FormData();
    formData.append('uploadId', uploadId);
    formData.append('basePath', params.path);
    formData.append('machineId', String(params.machineId));
    formData.append('authCertName', params.authCertName);
    formData.append('protocol', String(params.protocol));

    // 添加所有文件
    const paths: string[] = [];
    for (let i = 0; i < files.length; i++) {
        const file = files[i];
        const relativePath = (file as any).webkitRelativePath || file.name;
        formData.append('files', file);
        paths.push(relativePath);
    }

    // 添加路径数组
    paths.forEach((path) => {
        formData.append('paths', path);
    });

    // 使用 Api.upload 发起请求
    const { abort } = machineApi.uploadFolder.upload(formData, options);

    return { uploadId, abort };
}
