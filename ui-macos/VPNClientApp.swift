import SwiftUI

struct VPNClientApp: App {
    @StateObject private var viewModel = VPNViewModel()
    
    var body: some Scene {
        WindowGroup {
            ContentView()
                .environmentObject(viewModel)
        }
    }
}

struct ContentView: View {
    @EnvironmentObject var viewModel: VPNViewModel
    @State private var selectedCountry = ""
    
    var body: some View {
        VStack {
            Picker("Select Country", selection: $selectedCountry) {
                ForEach(viewModel.countries, id: \.self) { Text($0) }
            }
            Button("Connect") {
                viewModel.connectTo(selectedCountry)
            }
        }
        .onAppear {
            viewModel.fetchCountries()
        }
    }
}