package models

// SwapperManifest represents the root structure of the dlss-swapper manifest.json
type SwapperManifest struct {
	Version   string      `json:"version"`
	KnownDLLs []string    `json:"KnownDLLs"`
	DLSS      []DLLRecord `json:"DLSS"`
	DLSSG     []DLLRecord `json:"DLSS_G"`
	DLSSD     []DLLRecord `json:"DLSS_D"`
	FSR31DX12 []DLLRecord `json:"FSR_31_DX12"`
	FSR31VK   []DLLRecord `json:"FSR_31_VK"`
	XeSS      []DLLRecord `json:"XeSS"`
	XeSSFG    []DLLRecord `json:"XeSS_FG"`
}

// DLLRecord represents a single DLL version available for download
type DLLRecord struct {
	Version          string `json:"version"`
	VersionNumber    uint64 `json:"version_number"`
	InternalName     string `json:"internal_name,omitempty"`
	AdditionalLabel  string `json:"additional_label,omitempty"`
	MD5Hash          string `json:"md5_hash"`
	ZipMD5Hash       string `json:"zip_md5_hash,omitempty"`
	DownloadUrl      string `json:"download_url"`
	FileDescription  string `json:"file_description"`
	SignedDateTime   string `json:"signed_datetime"`
	IsSignatureValid bool   `json:"is_signature_valid"`
	IsDevFile        bool   `json:"is_dev_file"`
	FileSize         int64  `json:"file_size"`
	ZipFileSize      int64  `json:"zip_file_size"`

	// Local fields (not from JSON)
	Downloaded bool   `json:"downloaded"`
	LocalPath  string `json:"local_path,omitempty"`
	Type       string `json:"type"` // e.g. "dlss", "fsr", "xess"
}
