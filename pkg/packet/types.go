package packet

type ResourcePackStatus int

const (
	ResourcePackSuccessfullyLoaded ResourcePackStatus = iota
	ResourcePackDeclined
	ResourcePackFailedDownload
	ResourcePackAccepted
	ResourcePackDownloaded
	ResourcePackInvalidURL
	ResourcePackFailedReload
	ResourcePackDiscarded
)

func (s ResourcePackStatus) Validate() bool {
	return s >= ResourcePackSuccessfullyLoaded && s <= ResourcePackDiscarded
}

func (s ResourcePackStatus) String() string {
	switch s {
	case ResourcePackSuccessfullyLoaded:
		return "successfully_loaded"
	case ResourcePackDeclined:
		return "declined"
	case ResourcePackFailedDownload:
		return "failed_download"
	case ResourcePackAccepted:
		return "accepted"
	case ResourcePackDownloaded:
		return "downloaded"
	case ResourcePackInvalidURL:
		return "invalid_url"
	case ResourcePackFailedReload:
		return "failed_reload"
	case ResourcePackDiscarded:
		return "discarded"
	}
	return "unknown"
}
