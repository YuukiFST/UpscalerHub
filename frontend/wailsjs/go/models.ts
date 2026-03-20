export namespace main {
	
	export class AppUpdateInfo {
	    available: boolean;
	    version: string;
	    notes: string;
	    url: string;
	
	    static createFrom(source: any = {}) {
	        return new AppUpdateInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.available = source["available"];
	        this.version = source["version"];
	        this.notes = source["notes"];
	        this.url = source["url"];
	    }
	}
	export class ComponentStatus {
	    local: models.ComponentVersions;
	    remote: models.ComponentVersions;
	    optiScalerVersions: string[];
	    updateAvailable: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ComponentStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.local = this.convertValues(source["local"], models.ComponentVersions);
	        this.remote = this.convertValues(source["remote"], models.ComponentVersions);
	        this.optiScalerVersions = source["optiScalerVersions"];
	        this.updateAvailable = source["updateAvailable"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class GpuResult {
	    display: string;
	    info?: services.GpuInfo;
	
	    static createFrom(source: any = {}) {
	        return new GpuResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.display = source["display"];
	        this.info = this.convertValues(source["info"], services.GpuInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class InstallResult {
	    success: boolean;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new InstallResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	    }
	}

}

export namespace models {
	
	export class ScanExclusion {
	    Name: string;
	    PathSegment: string;
	
	    static createFrom(source: any = {}) {
	        return new ScanExclusion(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.PathSegment = source["PathSegment"];
	    }
	}
	export class RepositoryConfig {
	    RepoOwner: string;
	    RepoName: string;
	
	    static createFrom(source: any = {}) {
	        return new RepositoryConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.RepoOwner = source["RepoOwner"];
	        this.RepoName = source["RepoName"];
	    }
	}
	export class AppConfiguration {
	    App: RepositoryConfig;
	    OptiScaler: RepositoryConfig;
	    Fakenvapi: RepositoryConfig;
	    NukemFG: RepositoryConfig;
	    Language?: string;
	    ScanExclusions?: ScanExclusion[];
	
	    static createFrom(source: any = {}) {
	        return new AppConfiguration(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.App = this.convertValues(source["App"], RepositoryConfig);
	        this.OptiScaler = this.convertValues(source["OptiScaler"], RepositoryConfig);
	        this.Fakenvapi = this.convertValues(source["Fakenvapi"], RepositoryConfig);
	        this.NukemFG = this.convertValues(source["NukemFG"], RepositoryConfig);
	        this.Language = source["Language"];
	        this.ScanExclusions = this.convertValues(source["ScanExclusions"], ScanExclusion);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ComponentVersions {
	    OptiScalerVersion?: string;
	    FakenvapiVersion?: string;
	    NukemFGVersion?: string;
	
	    static createFrom(source: any = {}) {
	        return new ComponentVersions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.OptiScalerVersion = source["OptiScalerVersion"];
	        this.FakenvapiVersion = source["FakenvapiVersion"];
	        this.NukemFGVersion = source["NukemFGVersion"];
	    }
	}
	export class DLLRecord {
	    version: string;
	    version_number: number;
	    internal_name?: string;
	    additional_label?: string;
	    md5_hash: string;
	    zip_md5_hash?: string;
	    download_url: string;
	    file_description: string;
	    signed_datetime: string;
	    is_signature_valid: boolean;
	    is_dev_file: boolean;
	    file_size: number;
	    zip_file_size: number;
	    downloaded: boolean;
	    local_path?: string;
	    type: string;
	
	    static createFrom(source: any = {}) {
	        return new DLLRecord(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.version_number = source["version_number"];
	        this.internal_name = source["internal_name"];
	        this.additional_label = source["additional_label"];
	        this.md5_hash = source["md5_hash"];
	        this.zip_md5_hash = source["zip_md5_hash"];
	        this.download_url = source["download_url"];
	        this.file_description = source["file_description"];
	        this.signed_datetime = source["signed_datetime"];
	        this.is_signature_valid = source["is_signature_valid"];
	        this.is_dev_file = source["is_dev_file"];
	        this.file_size = source["file_size"];
	        this.zip_file_size = source["zip_file_size"];
	        this.downloaded = source["downloaded"];
	        this.local_path = source["local_path"];
	        this.type = source["type"];
	    }
	}
	export class Game {
	    name: string;
	    installPath: string;
	    platform: string;
	    appId: string;
	    executablePath: string;
	    coverImageUrl?: string;
	    isFavorite: boolean;
	    dlssVersion?: string;
	    dlssPath?: string;
	    dlssFrameGenVersion?: string;
	    dlssFrameGenPath?: string;
	    fsrVersion?: string;
	    fsrPath?: string;
	    xessVersion?: string;
	    xessPath?: string;
	    isOptiScalerInstalled: boolean;
	    optiScalerVersion?: string;
	
	    static createFrom(source: any = {}) {
	        return new Game(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.installPath = source["installPath"];
	        this.platform = source["platform"];
	        this.appId = source["appId"];
	        this.executablePath = source["executablePath"];
	        this.coverImageUrl = source["coverImageUrl"];
	        this.isFavorite = source["isFavorite"];
	        this.dlssVersion = source["dlssVersion"];
	        this.dlssPath = source["dlssPath"];
	        this.dlssFrameGenVersion = source["dlssFrameGenVersion"];
	        this.dlssFrameGenPath = source["dlssFrameGenPath"];
	        this.fsrVersion = source["fsrVersion"];
	        this.fsrPath = source["fsrPath"];
	        this.xessVersion = source["xessVersion"];
	        this.xessPath = source["xessPath"];
	        this.isOptiScalerInstalled = source["isOptiScalerInstalled"];
	        this.optiScalerVersion = source["optiScalerVersion"];
	    }
	}
	
	
	export class SwapperManifest {
	    version: string;
	    KnownDLLs: string[];
	    DLSS: DLLRecord[];
	    DLSS_G: DLLRecord[];
	    DLSS_D: DLLRecord[];
	    FSR_31_DX12: DLLRecord[];
	    FSR_31_VK: DLLRecord[];
	    XeSS: DLLRecord[];
	    XeSS_FG: DLLRecord[];
	
	    static createFrom(source: any = {}) {
	        return new SwapperManifest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.KnownDLLs = source["KnownDLLs"];
	        this.DLSS = this.convertValues(source["DLSS"], DLLRecord);
	        this.DLSS_G = this.convertValues(source["DLSS_G"], DLLRecord);
	        this.DLSS_D = this.convertValues(source["DLSS_D"], DLLRecord);
	        this.FSR_31_DX12 = this.convertValues(source["FSR_31_DX12"], DLLRecord);
	        this.FSR_31_VK = this.convertValues(source["FSR_31_VK"], DLLRecord);
	        this.XeSS = this.convertValues(source["XeSS"], DLLRecord);
	        this.XeSS_FG = this.convertValues(source["XeSS_FG"], DLLRecord);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace nvapi {
	
	export class PresetResult {
	    success: boolean;
	    error?: string;
	    dlssPreset?: number;
	    dlssdPreset?: number;
	    foundProfile?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new PresetResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.error = source["error"];
	        this.dlssPreset = source["dlssPreset"];
	        this.dlssdPreset = source["dlssdPreset"];
	        this.foundProfile = source["foundProfile"];
	    }
	}

}

export namespace services {
	
	export class GpuInfo {
	    name: string;
	    vendor: string;
	    driverVersion: string;
	    videoMemoryMB: number;
	    icon: string;
	
	    static createFrom(source: any = {}) {
	        return new GpuInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.vendor = source["vendor"];
	        this.driverVersion = source["driverVersion"];
	        this.videoMemoryMB = source["videoMemoryMB"];
	        this.icon = source["icon"];
	    }
	}

}

