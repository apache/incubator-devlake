const jiraUtil = {
  calculateLeadTime: (startDate, endDate) => {
		try {
			let seconds = endDate.getTime() - startDate.getTime() / 1000
			return new Date(seconds * 1000).toISOString().substr(11, 8)
		} catch (error) {
			throw new Error('Calculate Lead time failed')
		}
  },
};

module.exports = jiraUtil;
