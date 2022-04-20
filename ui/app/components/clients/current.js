import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { isAfter, startOfMonth } from 'date-fns';
import { action } from '@ember/object';
export default class Current extends Component {
  chartLegend = [
    { key: 'entity_clients', label: 'entity clients' },
    { key: 'non_entity_clients', label: 'non-entity clients' },
  ];
  @tracked selectedNamespace = null;
  @tracked namespaceArray = this.byNamespaceTotalClients.map((namespace) => {
    return { name: namespace['label'], id: namespace['label'] };
  });

  @tracked selectedAuthMethod = null;
  @tracked authMethodOptions = [];

  get latestUpgradeData() {
    // e.g. {id: '1.9.0', previousVersion: null, timestampInstalled: '2021-11-03T10:23:16Z'}
    // version id is 1.9.0 or earliest upgrade post 1.9.0, timestamp is RFC3339
    // need to return the last in the array as the data is in ascending order
    return this.args.model.versionHistory.at(-1) || null;
  }

  // Response total client count data by namespace for current/partial month
  get byNamespaceTotalClients() {
    return this.args.model.monthly?.byNamespaceTotalClients || [];
  }

  // Response new client count data by namespace for current/partial month
  get byNamespaceNewClients() {
    return this.args.model.monthly?.byNamespaceNewClients || [];
  }

  get isGatheringData() {
    // return true if tracking IS enabled but no data collected yet
    return (
      this.args.model.config?.enabled === 'On' &&
      this.byNamespaceTotalClients.length === 0 &&
      this.byNamespaceNewClients.length === 0
    );
  }

  get hasAttributionData() {
    if (this.selectedAuthMethod) return false;
    if (this.selectedNamespace) {
      return this.authMethodOptions.length > 0;
    }
    return this.totalUsageCounts.clients !== 0 && !!this.totalClientsData;
  }

  get filteredTotalData() {
    const namespace = this.selectedNamespace;
    const auth = this.selectedAuthMethod;
    if (!namespace && !auth) {
      return this.byNamespaceTotalClients;
    }
    if (!auth) {
      return this.byNamespaceTotalClients.find((ns) => ns.label === namespace);
    }
    return this.byNamespaceTotalClients
      .find((ns) => ns.label === namespace)
      .mounts?.find((mount) => mount.label === auth);
  }

  get filteredNewData() {
    const namespace = this.selectedNamespace;
    const auth = this.selectedAuthMethod;
    if (!namespace && !auth) {
      return this.byNamespaceNewClients;
    }
    if (!auth) {
      return this.byNamespaceNewClients.find((ns) => ns.label === namespace);
    }
    return this.byNamespaceNewClients
      .find((ns) => ns.label === namespace)
      .mounts?.find((mount) => mount.label === auth);
  }

  get versionUpdateText() {
    if (!this.latestUpgradeData) {
      return false;
    }

    let versionDate = new Date(this.latestUpgradeData.timestampInstalled);
    if (isAfter(versionDate, startOfMonth(new Date()))) {
      if (this.minorReleaseNumber === 9) {
        return `Vault was upgraded to 
        ${this.latestUpgradeData.id} this month. How we count clients changed in 1.9, so keep that in mind when looking at the data below. `;
      } else if (this.minorReleaseNumber === 10) {
        return `Vault was upgraded to 
        ${this.latestUpgradeData.id} this month. We added new client breakdowns for 1.10 attribution data only, so keep that in mind when looking at the data below. `;
      } else {
        return null;
      }
    } else {
      return null;
    }
  }

  get minorReleaseNumber() {
    let upgradeVersion = this.latestUpgradeData.id;
    let upgradeVersionArray = upgradeVersion.split('.');
    return Number(upgradeVersionArray[1]);
  }

  // top level TOTAL client counts for current/partial month
  get totalUsageCounts() {
    return this.selectedNamespace ? this.filteredTotalData : this.args.model.monthly?.total;
  }

  get newUsageCounts() {
    return this.selectedNamespace ? this.filteredNewData : this.args.model.monthly?.new;
  }

  // total client data for horizontal bar chart in attribution component
  get totalClientsData() {
    if (this.selectedNamespace) {
      return this.filteredTotalData?.mounts || null;
    } else {
      return this.byNamespaceTotalClients;
    }
  }

  // new client data for horizontal bar chart in attribution component
  get newClientsData() {
    if (this.selectedNamespace) {
      return this.filteredNewData?.mounts || null;
    } else {
      return this.byNamespaceNewClients;
    }
  }

  get responseTimestamp() {
    return this.args.model.monthly?.responseTimestamp;
  }

  // ACTIONS
  @action
  selectNamespace([value]) {
    // value comes in as [namespace0]
    this.selectedNamespace = value;
    if (!value) {
      this.authMethodOptions = [];
      // on clear, also make sure auth method is cleared
      this.selectedAuthMethod = null;
    } else {
      // Side effect: set auth namespaces
      const mounts = this.filteredTotalData.mounts?.map((mount) => ({
        id: mount.label,
        name: mount.label,
      }));
      this.authMethodOptions = mounts;
    }
  }

  @action
  setAuthMethod([authMount]) {
    this.selectedAuthMethod = authMount;
  }
}
