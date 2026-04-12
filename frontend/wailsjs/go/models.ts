export namespace database {
	
	export class Data {
	    Mods: models.Mod[];
	    Profiles: models.Profile[];
	    ModProfiles: models.ModProfile[];
	
	    static createFrom(source: any = {}) {
	        return new Data(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Mods = this.convertValues(source["Mods"], models.Mod);
	        this.Profiles = this.convertValues(source["Profiles"], models.Profile);
	        this.ModProfiles = this.convertValues(source["ModProfiles"], models.ModProfile);
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

export namespace models {
	
	export class Mod {
	    Dir: string;
	    Name: string;
	    Category: string;
	    Active: boolean;
	    InProfile: boolean;
	    // Go type: time
	    LastModified: any;
	
	    static createFrom(source: any = {}) {
	        return new Mod(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Dir = source["Dir"];
	        this.Name = source["Name"];
	        this.Category = source["Category"];
	        this.Active = source["Active"];
	        this.InProfile = source["InProfile"];
	        this.LastModified = this.convertValues(source["LastModified"], null);
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
	export class ModProfile {
	    ModDir: string;
	    ProfileID: number;
	
	    static createFrom(source: any = {}) {
	        return new ModProfile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ModDir = source["ModDir"];
	        this.ProfileID = source["ProfileID"];
	    }
	}
	export class Profile {
	    Id: number;
	    Name: string;
	    Path: string;
	    Category: string;
	    Active: boolean;
	    AutoCreated: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Profile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.Name = source["Name"];
	        this.Path = source["Path"];
	        this.Category = source["Category"];
	        this.Active = source["Active"];
	        this.AutoCreated = source["AutoCreated"];
	    }
	}
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

