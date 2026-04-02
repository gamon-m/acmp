export namespace main {
	
	export class Settings {
	    assetto_corsa_path: string;
	    mods_path: string;
	    automatic_profiles: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.assetto_corsa_path = source["assetto_corsa_path"];
	        this.mods_path = source["mods_path"];
	        this.automatic_profiles = source["automatic_profiles"];
	    }
	}

}

