public partial class MainWindow : Window {
    private readonly HttpClient _client = new HttpClient();
    private ObservableCollection<string> _countries = new ObservableCollection<string>();
    
    public MainWindow() {
        InitializeComponent();
        CountryComboBox.ItemsSource = _countries;
        LoadCountries();
    }

    private async void LoadCountries() {
        var response = await _client.GetAsync("http://localhost:8080/servers/countries");
        if (response.IsSuccessStatusCode) {
            var content = await response.Content.ReadAsStringAsync();
            var countries = JsonConvert.DeserializeObject<List<string>>(content);
            _countries.Clear();
            foreach (var country in countries) {
                _countries.Add(country);
            }
        }
    }

    private void ConnectButton_Click(object sender, RoutedEventArgs e) {
        var selectedCountry = CountryComboBox.SelectedItem?.ToString();
        // Implement connection logic
    }
}