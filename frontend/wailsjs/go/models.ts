export namespace models {
	
	export class APIKeyDTO {
	    id: number;
	    keyName: string;
	    keyValue: string;
	    url: string;
	    expiryDate: string;
	    registeredDate: string;
	    memo: string;
	    createdAt: string;
	    updatedAt: string;
	    expiryStatus: number;
	
	    static createFrom(source: any = {}) {
	        return new APIKeyDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.keyName = source["keyName"];
	        this.keyValue = source["keyValue"];
	        this.url = source["url"];
	        this.expiryDate = source["expiryDate"];
	        this.registeredDate = source["registeredDate"];
	        this.memo = source["memo"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	        this.expiryStatus = source["expiryStatus"];
	    }
	}
	export class CreateKeyRequest {
	    keyName: string;
	    keyValue: string;
	    url: string;
	    expiryDate: string;
	    registeredDate: string;
	    memo: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateKeyRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.keyName = source["keyName"];
	        this.keyValue = source["keyValue"];
	        this.url = source["url"];
	        this.expiryDate = source["expiryDate"];
	        this.registeredDate = source["registeredDate"];
	        this.memo = source["memo"];
	    }
	}
	export class KeyFilter {
	    keyName: string;
	    dateFrom: string;
	    dateTo: string;
	
	    static createFrom(source: any = {}) {
	        return new KeyFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.keyName = source["keyName"];
	        this.dateFrom = source["dateFrom"];
	        this.dateTo = source["dateTo"];
	    }
	}
	export class UpdateKeyRequest {
	    id: number;
	    keyName: string;
	    keyValue: string;
	    url: string;
	    expiryDate: string;
	    registeredDate: string;
	    memo: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateKeyRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.keyName = source["keyName"];
	        this.keyValue = source["keyValue"];
	        this.url = source["url"];
	        this.expiryDate = source["expiryDate"];
	        this.registeredDate = source["registeredDate"];
	        this.memo = source["memo"];
	    }
	}

}

