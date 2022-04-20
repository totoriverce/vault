import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { isSameMonth, isAfter, isBefore } from 'date-fns';
import getStorage from 'vault/lib/token-storage';
import { ARRAY_OF_MONTHS } from 'core/utils/date-formatters';

const INPUTTED_START_DATE = 'vault:ui-inputted-start-date';

export default class History extends Component {
  @service store;
  @service version;

  arrayOfMonths = ARRAY_OF_MONTHS;

  chartLegend = [
    { key: 'entity_clients', label: 'entity clients' },
    { key: 'non_entity_clients', label: 'non-entity clients' },
  ];

  // FOR START DATE EDIT & MODAL //
  months = Array.from({ length: 12 }, (item, i) => {
    return new Date(0, i).toLocaleString('en-US', { month: 'long' });
  });
  years = Array.from({ length: 5 }, (item, i) => {
    return new Date().getFullYear() - i;
  });
  currentDate = new Date();
  currentYear = this.currentDate.getFullYear(); // integer of year
  currentMonth = this.currentDate.getMonth(); // index of month

  @tracked isEditStartMonthOpen = false;
  @tracked startMonth = null;
  @tracked startYear = null;
  @tracked allowedMonthMax = 12;
  @tracked disabledYear = null;

  // FOR HISTORY COMPONENT //

  // RESPONSE
  @tracked endTimeFromResponse = this.args.model.endTimeFromResponse;
  @tracked startTimeFromResponse = this.args.model.startTimeFromLicense; // ex: ['2021', 3] is April 2021 (0 indexed)
  @tracked startTimeRequested = null;
  @tracked queriedActivityResponse = null;

  // SEARCH SELECT
  @tracked selectedNamespace = null;
  @tracked namespaceArray = this.getActivityResponse.byNamespace
    ? this.getActivityResponse.byNamespace.map((namespace) => ({
        name: namespace.label,
        id: namespace.label,
      }))
    : [];
  @tracked selectedAuthMethod = null;
  @tracked authMethodOptions = [];

  // TEMPLATE MESSAGING
  @tracked noActivityDate = '';
  @tracked responseRangeDiffMessage = null;
  @tracked isLoadingQuery = false;
  @tracked licenseStartIsCurrentMonth = this.args.model.activity?.isLicenseDateError || false;
  @tracked errorObject = null;

  get versionText() {
    return this.version.isEnterprise
      ? {
          label: 'Billing start month',
          description:
            'This date comes from your license, and defines when client counting starts. Without this starting point, the data shown is not reliable.',
          title: 'No billing start date found',
          message:
            'In order to get the most from this data, please enter your billing period start month. This will ensure that the resulting data is accurate.',
        }
      : {
          label: 'Client counting start date',
          description:
            'This date is when client counting starts. Without this starting point, the data shown is not reliable.',
          title: 'No start date found',
          message:
            'In order to get the most from this data, please enter a start month above. Vault will calculate new clients starting from that month.',
        };
  }

  get isDateRange() {
    return !isSameMonth(
      new Date(this.getActivityResponse.startTime),
      new Date(this.getActivityResponse.endTime)
    );
  }

  get latestUpgradeData() {
    // {id: '1.9.0', previousVersion: null, timestampInstalled: '2021-11-03T10:23:16Z'}
    // version id is 1.9.0 or earliest upgrade post 1.9.0, timestamp is RFC3339
    // need to return the last in the array as the data is in ascending order
    return this.args.model.versionHistory.at(-1) || null;
  }

  get allMinorVersionsInRange() {
    let newArray = [];
    // ARG TODO testing unsure if the getActivityResponse is the correct one here.
    let startTime = new Date(this.getActivityResponse.startTime);
    let endTime = new Date(this.getActivityResponse.endTime);
    this.args.model.versionHistory.forEach((object) => {
      let versionDate = new Date(object.timestampInstalled);
      // compare against the API response to see if date is within the range of returned response.
      if (isAfter(versionDate, startTime) && isBefore(versionDate, endTime)) {
        // push the minor version number e.g. 9
        let upgradeVersionArray = object.id.split('.');
        newArray.push(upgradeVersionArray[1]);
      }
    });
    return newArray; // return [9,9,10] shows all the minor numbers that fall within the range of the returned API response.
  }

  get versionUpdateText() {
    // count how many times a minor number occurs inside the array of minors within the API range.
    let occurrences = this.allMinorVersionsInRange.reduce(function (acc, curr) {
      return acc[curr] ? ++acc[curr] : (acc[curr] = 1), acc;
    }, {});
    // if count of minor 9 is greater than one and there are no occurrences of minor 10, show the minor 9 specific warning
    if (occurrences[9] >= 1 && !occurrences[10]) {
      return `Vault was upgraded to ${this.latestUpgradeData.id} during this time range. How we count clients changed in 1.9, so keep that in mind when looking at the data below. `;
    } else if (occurrences[10] >= 1 && !occurrences[9]) {
      // show the minor 10 specific warning
      return `Vault was upgraded to ${this.latestUpgradeData.id} during this time range. We added monthly breakdowns for 1.10 data only, so keep that in mind when looking at the data below. `;
    } else if (occurrences[10] >= 1 && occurrences[9] >= 1) {
      // you have both 9 and 10 as minors. Show generic warning.
      return `Vault was upgraded to 1.9 and 1.10 during this time range. How we count clients changed in 1.9 and we added monthly breakdowns for 1.10 data only. Keep this in mind when looking at the data below. `;
    }
    // null if there are no 9 or 10 minor versions within the range returned by the API
    return null;
  }

  get startTimeDisplay() {
    if (!this.startTimeFromResponse) {
      return null;
    }
    let month = this.startTimeFromResponse[1];
    let year = this.startTimeFromResponse[0];
    return `${this.arrayOfMonths[month]} ${year}`;
  }

  get endTimeDisplay() {
    if (!this.endTimeFromResponse) {
      return null;
    }
    let month = this.endTimeFromResponse[1];
    let year = this.endTimeFromResponse[0];
    return `${this.arrayOfMonths[month]} ${year}`;
  }

  // GETTERS FOR RESPONSE & DATA

  // on init API response uses license start_date, getter updates when user queries dates
  get getActivityResponse() {
    return this.queriedActivityResponse || this.args.model.activity;
  }

  get hasAttributionData() {
    if (this.selectedAuthMethod) return false;
    if (this.selectedNamespace) {
      return this.authMethodOptions.length > 0;
    }
    return !!this.totalClientsData && this.totalUsageCounts && this.totalUsageCounts.clients !== 0;
  }

  // top level TOTAL client counts for given date range
  get totalUsageCounts() {
    return this.selectedNamespace ? this.filteredActivity : this.getActivityResponse.total;
  }

  // total client data for horizontal bar chart in attribution component
  get totalClientsData() {
    if (this.selectedNamespace) {
      return this.filteredActivity?.mounts || null;
    } else {
      return this.getActivityResponse?.byNamespace;
    }
  }

  get responseTimestamp() {
    return this.getActivityResponse.responseTimestamp;
  }

  get byMonthTotalClients() {
    return this.getActivityResponse?.byMonth;
  }

  get byMonthNewClients() {
    return this.byMonthTotalClients.map((m) => m.new_clients);
  }

  get filteredActivity() {
    const namespace = this.selectedNamespace;
    const auth = this.selectedAuthMethod;
    if (!namespace && !auth) {
      return this.getActivityResponse;
    }
    if (!auth) {
      return this.getActivityResponse.byNamespace.find((ns) => ns.label === namespace);
    }
    return this.getActivityResponse.byNamespace
      .find((ns) => ns.label === namespace)
      .mounts?.find((mount) => mount.label === auth);
  }

  @action
  async handleClientActivityQuery(month, year, dateType) {
    this.isEditStartMonthOpen = false;
    if (dateType === 'cancel') {
      return;
    }
    // clicked "Current Billing period" in the calendar widget
    if (dateType === 'reset') {
      this.startTimeRequested = this.args.model.startTimeFromLicense;
      this.endTimeRequested = null;
    }
    // clicked "Edit" Billing start month in History which opens a modal.
    if (dateType === 'startTime') {
      let monthIndex = this.arrayOfMonths.indexOf(month);
      this.startTimeRequested = [year.toString(), monthIndex]; // ['2021', 0] (e.g. January 2021)
      this.endTimeRequested = null;
    }
    // clicked "Custom End Month" from the calendar-widget
    if (dateType === 'endTime') {
      // use the currently selected startTime for your startTimeRequested.
      this.startTimeRequested = this.startTimeFromResponse;
      this.endTimeRequested = [year.toString(), month]; // endTime comes in as a number/index whereas startTime comes in as a month name. Hence the difference between monthIndex and month.
    }

    try {
      this.isLoadingQuery = true;
      let response = await this.store.queryRecord('clients/activity', {
        start_time: this.startTimeRequested,
        end_time: this.endTimeRequested,
      });
      if (response.id === 'no-data') {
        // empty response (204) is the only time we want to update the displayed date with the requested time
        this.startTimeFromResponse = this.startTimeRequested;
        this.noActivityDate = this.startTimeDisplay;
      } else {
        // note: this.startTimeDisplay (getter) is updated by the @tracked startTimeFromResponse
        this.startTimeFromResponse = response.formattedStartTime;
        this.endTimeFromResponse = response.formattedEndTime;
        this.storage().setItem(INPUTTED_START_DATE, this.startTimeFromResponse);
      }
      this.queriedActivityResponse = response;
      this.licenseStartIsCurrentMonth = response.isLicenseDateError;
      // compare if the response startTime comes after the requested startTime. If true throw a warning.
      // only display if they selected a startTime
      if (
        dateType === 'startTime' &&
        isAfter(
          new Date(this.getActivityResponse.startTime),
          new Date(this.startTimeRequested[0], this.startTimeRequested[1])
        )
      ) {
        this.responseRangeDiffMessage = `You requested data from ${month} ${year}. We only have data from ${this.startTimeDisplay}, and that is what is being shown here.`;
      } else {
        this.responseRangeDiffMessage = null;
      }
    } catch (e) {
      this.errorObject = e;
      return e;
    } finally {
      this.isLoadingQuery = false;
    }
  }

  @action
  handleCurrentBillingPeriod() {
    this.handleClientActivityQuery(0, 0, 'reset');
  }

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
      const mounts = this.filteredActivity.mounts?.map((mount) => ({
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

  // FOR START DATE MODAL
  @action
  selectStartMonth(month, event) {
    this.startMonth = month;
    // disables months if in the future
    this.disabledYear = this.months.indexOf(month) >= this.currentMonth ? this.currentYear : null;
    event.close();
  }

  @action
  selectStartYear(year, event) {
    this.startYear = year;
    this.allowedMonthMax = year === this.currentYear ? this.currentMonth : 12;
    event.close();
  }

  storage() {
    return getStorage();
  }
}
